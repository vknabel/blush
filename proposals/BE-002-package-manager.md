# Package Manager and Cavefile

- **Proposal:** BE-002
- **Status:** Draft
- **Authors:** [@vknabel](https://github.com/vknabel), [@blushling](https://github.com/blushling)

## Overview

Blush ships with a Go-based package manager that coordinates dependency
resolution across one or more registries. `PackageManager` stores the available
registries and exposes an `Install` helper that prepares an `InstallationTask`
for a parsed `Cavefile`. The default constructor wires in a Git-backed registry
rooted at `git/` inside the provided filesystem, which mirrors how the runtime
stores cached packages.

## Cavefile data model

The package manager consumes a simplified representation of the manifest:

```go
type Cavefile struct {
        Dependencies []Dependency
}

type Dependency struct {
        ImportName string
        Source     string
        Predicate  version.Predicate
}
```

`ImportName` records the identifier a project expects to use for an imported
package, `Source` points to the upstream location (such as a Git URL), and
`Predicate` captures the semantic-version constraint that must be satisfied by a
resolved release.

Blush manifests can contain additional declarations for automation or tooling.
The package manager, however, only inspects the dependency list. For example,
the sample project declares both dependency metadata and a task definition, but
only the entries under `@cave.Dependencies` affect installation.

### Example Cavefile

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

## Installation workflow

Executing `InstallationTask.Run` performs dependency resolution in three phases:

1. Initialise the work queue from the manifest if this is the first run.
2. Ask every configured registry for locally cached packages and index the
   results by their source string.
3. For each dependency, prefer a matching local package. If none are available,
   query each registry for remote versions that satisfy the predicate, resolve
   (clone) the first match, and append it to the completed list.

If no registry can serve a dependency, the task aborts with an error. Completed
packages remain attached to the task instance so tooling can reuse the results in
subsequent steps.

## Registry abstraction

Registries implement a common interface that supports two operations: enumerate
locally resolved packages and discover remote versions that meet a set of
predicates. The API also defines the `Package`, `ResolvedPackage`, and
`ResolvedModule` contracts used throughout the installer. Blush expects packages
to live under a `git/<package>/<version>/` hierarchy beneath either the standard
library or user-controlled cache directories.

### Git registry provider

The built-in `GitRegistry` keeps its workspace inside the filesystem supplied to
`PackageManager.New`. During local discovery it traverses each repository folder,
opens stored clones, and exposes tag aliases that refer to cached versions. When
remote discovery is requested, it lists tags from the upstream repository,
filters them against the predicates, sorts versions in descending order, and
returns descriptors that can be cloned on demand. Cloning writes into a
`<source>/<version>/` subdirectory whose name is mangled to remain filesystem
friendly before returning a `localGitPackage`.

Resolved Git packages delegate module discovery to the filesystem-based helper in
`registry/fsmodule`. This helper walks the repository for `.blush` sources,
groups them by directory, and exposes each group as a module accessible via a
logical URI derived from the package source and path.

## Limitations and open questions

The current implementation leaves several areas for future work:

- Transitive dependency installation is unimplemented; the task processes only
  the top-level entries from the manifest queue.
- `Dependency.ImportName` is recorded but not yet consulted during resolution,
  so package aliases do not influence how modules are installed.
- Only the Git registry backend is bundled today. Additional providers are
  required to support alternative transports or storage layouts.
- When multiple remote versions satisfy the predicate, the installer currently
  resolves the first candidate returned by the registry after sorting, leaving
  room for richer selection strategies.

Documenting these constraints clarifies the baseline behaviour and helps guide
future extensions to Blush's package ecosystem.
