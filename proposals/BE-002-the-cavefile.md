# The Cavefile

- **Proposal:** BE-002
- **Status:** In Progress
- **Author:** [@vknabel](https://github.com/vknabel), [@blushling](https://github.com/blushling)

## Introduction

This proposal describes the current Blush package manager and the `Cavefile` manifest that drives it. The goal is to document the partially implemented behaviour so users and contributors understand how dependencies are declared, resolved, and cached today.

## Motivation

Blush projects will rely on external modules for language extensions and tooling.
Capturing the present-day and near future behaviour of the package manager clarifies how packages are discovered, how versions are selected, and how registries interact with the filesystem cache. Documenting the `Cavefile` structure likewise gives
users a reference for authoring manifests that match the implementation.
## Proposed Solution

Packages declare dependencies in a `Cavefile` module that Blush parses at build or install time. The manifest itself will not be executed. The package manager solely works on the type system.

The package manager itself caches Git repositories on disk and coordinates one or more registries to provide the requested packages. Currently Git, local files and built-in registries are supported.

As of now the package manager does not traverse transitive dependencies, leaving that for a future enhancement.

### Cavefile Dependencies

The data structure with the `@cave.Dependencies()` annotation will be used to declare dependencies in a `Cavefile`.

```blush
import cave

@cave.Dependencies()
data Dependencies {
  // Standard library packages are available by default
  @cave.Stdlib("prelude")
  prelude

  // References a local package. The path is relative to the Cavefile location.
  @cave.Local("../some-local-package")
  helpers // importable as helpers

  // References a package hosted in a Git repository.
  @cave.Git("https://github.com/vknabel/blush")
  @cave.Version(">0.1.0")
  future // importable as future
}
```

Only the entries produced by `@cave.Dependencies` participate in dependency
resolution, but additional declarations such as tasks can live alongside them.

### Cavefile Tasks

Additionally to dependencies, a `Cavefile` can declare tasks that the Blush CLI can execute.

```blush
@tasks.Name("generate")
@tasks.Help("Generates something")
@tasks.Exec("tasks/generate.blush")
data GenerateTask {
	@Bool
	@tasks.Flag()
  @tasks.Name("dry")
	@tasks.Help("If true, only simulates the generation")
	isDryRun

	@String
	@tasks.Arg()
	positional
}
```

This will declare a task named `generate` that can be executed with `blush generate`.
Upon execution, the `tasks/generate.blush` script will be run with the provided flags and arguments.

## Detailed Design

### Dependencies manifest structure

The `pkgmanager` will use the `Cavefile`, search for the `@cave.Dependencies()` data structure and parse its fields as follows:

- **Field name** – the import name for the dependency package. Submodules are supported by
  allowing dot notation, e.g. `foo.bar` imports the `bar` submodule from the
  `foo` package.
- **Source** - determined by the presence of `@cave.Git`, `@cave.Local`, or
  `@cave.Stdlib` annotations on the field.
- **Version predicate** – extracted from the `@cave.Version` annotation if
  present. If omitted, any version is acceptable.

### Package manager workflow

`pkgmanager.New` initialises a `PackageManager` with the configured registries.
The default constructor mounts a Git registry under the `git/` directory of the
provided filesystem so cached repositories are stored beneath that path.

`PackageManager.Install` creates an `InstallationTask` bound to a parsed
Cavefile. Running the task performs the following steps:

1. **Queue setup** – the first execution copies the manifest dependencies into a
   work queue so repeated runs reuse the same slice.
2. **Local discovery** – each registry reports the packages already cached on
   disk. Results are grouped by source so matching versions can be reused.
3. **Remote resolution** – unmet dependencies trigger `DiscoverPackageVersions`
   on every registry. The installer selects the first offered version (registries
   return results sorted from newest to oldest), resolves it to a local clone,
   and records the package.

If no registry can satisfy a dependency, the run terminates with an error. The
installer does not yet traverse transitive dependencies, leaving that work for a
future enhancement.

### Registry abstraction

Registries implement the `registry.Provider` interface. They surface:

- `Discover`, which returns all locally available packages and versions.
- `DiscoverPackageVersions`, which lists remote versions matching supplied
  predicates.

Packages resolved from either path expose module discovery helpers so the rest
of the toolchain can load `.blush` sources. Modules advertise their logical URI
and enumerate the files contained within the checkout.

### Git registry provider

The bundled `GitRegistry` manages repositories inside its root filesystem. Local
packages are stored as `git/<source>/<version>/` beneath the registry root.
Local discovery iterates those directories, opens each Git worktree, and lists
any tags that match the checked-out commit. Remote discovery connects to the
upstream repository, collects available tags, filters them by the provided
predicates, and sorts them in descending semantic-version order before returning
packages to the installer.

When a remote version is selected, the registry clones the tagged commit into a
mangled `<source>/<version>/` directory. Module enumeration delegates to the
filesystem module discovery helper so every `.blush` source within the repo is
published to the toolchain.

### Task execution and parsing

The `cave.tasks` package provides annotations and helpers to declare and execute tasks.
Data structures may be tasks when annotated with `@tasks.Exec` to execute files, `@tasks.Call` to call functions or `@tasks.Import` to reuse existing tasks.

They will be parsed by the CLI and registered as commands. By default the command name is the lowercased data name, but it may be overridden with `@tasks.Name`. A help text may be provided with `@tasks.Help`.

Tasks may declare flags and positional arguments by annotating fields with `@tasks.Flag` and `@tasks.Arg`.

## Changes to the Standard Library

Introduces the `cave` and `cave.tasks` modules to the standard library.

- `cave`:
  - `Dependencies` annotation
  - an enum for `Source` with values `Stdlib`, `Local`, and `Git`
  - `Stdlib`, `Local`, `Git`, and `Version` annotations
- `cave.tasks`:
  - an enum for `Task` with values `Exec`, `Call`, and `Import`
  - `Exec`, `Call`, and `Import` annotations for task declarations
  - `Name`, `Alias`, `Help`, `Short`, `Flag`, and `Arg` annotations for tasks and their fields

## Alternatives Considered

- Using a different manifest format such as JSON or YAML was considered, but
  Blush's strong typing and annotation system makes it straightforward to
  declare dependencies and tasks directly in Blush code.
- Implementing transitive dependency resolution was considered, but deferred to a future enhancement to keep the initial implementation simpler and focused on
  direct dependencies only.
- Supporting additional registry types (e.g., HTTP-based registries) was considered, but the initial implementation focuses on Git and local files to establish a solid foundation before expanding to other sources.
- Using a different approach for task declaration, such as a dedicated task configuration file, was considered, but integrating tasks into the `Cavefile` keeps related configurations together and leverages Blush's type system.
- A `package.blush` file was considered, but the `.blush` file extension would require additional tooling support and could lead to confusion with regular Blush source files.

## Acknowledgements

Thanks to the Blush maintainers for building the initial package manager and
Cavefile tooling that this document captures.
