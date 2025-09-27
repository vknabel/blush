# Package Manager and Cavefile

- **Proposal:** BE-002
- **Status:** Draft
- **Author:** [@vknabel](https://github.com/vknabel), [@blushling](https://github.com/blushling)

## Introduction

This proposal describes the current Blush package manager and the Cavefile
manifest that drives it. The goal is to document the behaviour already shipped
with Blush so users and contributors understand how dependencies are declared,
resolved, and cached today.

## Motivation

Blush projects rely on external modules for language extensions and tooling.
Capturing the present-day behaviour of the package manager clarifies how
packages are discovered, how versions are selected, and how registries interact
with the filesystem cache. Documenting the Cavefile structure likewise gives
users a reference for authoring manifests that match the implementation.

## Proposed Solution

Blush keeps the existing Go implementation as the reference package manager.
Projects declare dependencies in a `Cavefile` module that Blush parses at build
or install time. The package manager reads those dependency entries, prepares an
installation task, and coordinates one or more registries to provide the
requested packages. Registries expose both the locally cached packages and the
remote versions they can supply, allowing the installer to reuse previously
cloned repositories or fetch new ones as needed.

### Cavefile example

```blush
import cave
import cave.tasks

@cave.Dependencies()
data Dependencies {
        @cave.Prelude()
        data Prelude {
                @tasks.Use("test")
                test
        }
}

@tasks.Name("generate")
@tasks.Help("")
data GenerateTask {
        @Bool
        @tasks.Flag("dry")
        @tasks.Help("")
        isDryRun

        @String
        @tasks.Arg()
        positional
}
```

Only the entries produced by `@cave.Dependencies` participate in dependency
resolution, but additional declarations such as tasks can live alongside them.

## Detailed Design

### Cavefile manifest structure

The parser records each dependency as a value containing:

- **Import name** – the identifier made available to `import` statements.
- **Source** – the upstream location of the package, typically a Git URL.
- **Version predicate** – a semantic-version constraint the selected package
  must satisfy.

These fields mirror `cavefile.Dependency` and are consumed by the installer. The
package manager currently does not act on the import name during resolution but
retains it for future use.

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

## Changes to the Standard Library

None.

## Alternatives Considered

- Parsing the Cavefile into language-specific data structures other than the
  Go definitions used today. The current approach keeps the manifest simple and
  matches the implementation with minimal transformation.
- Allowing registries to decide how to order version candidates rather than
  sorting them centrally. Having the Git registry return versions in descending
  order ensures predictable selection even before richer resolution policies are
  introduced.

## Acknowledgements

Thanks to the Blush maintainers for building the initial package manager and
Cavefile tooling that this document captures.
