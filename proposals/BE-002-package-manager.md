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

The installer works with a high-level manifest called the Cavefile. When Blush
parses this document it records each dependency's

- **import name** – the identifier used in `import` statements,
- **source** – the upstream location (today this is normally a Git URL), and
- **version requirement** – a semantic-version predicate that must be satisfied
  by the resolved package.

Other declarations in the Cavefile continue to be available to author project
automation, but only the dependencies list participates in package resolution.
For example, the sample project below defines both dependency metadata and a
task, yet the installer only inspects the entries created by
`@cave.Dependencies`.

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

Creating an installation task (`PackageManager.Install`) snapshots the manifest
and defers the actual work until the task runs. When executed, the installer
performs three high-level steps:

1. **Prepare the queue.** The dependencies from the Cavefile seed the task the
   first time it runs so repeated invocations can reuse progress.
2. **Discover local packages.** Every configured registry is asked to describe
   packages that are already cached on disk. The results are grouped by source
   so the installer can quickly match them to dependencies.
3. **Resolve missing versions.** Dependencies that cannot be satisfied locally
   trigger remote discovery. Registries return the versions that satisfy the
   predicate, the installer picks the first candidate, clones it into the cache,
   and records the resolved package.

If no registry can satisfy a dependency, the task reports an error and stops.
Successful resolutions remain attached to the task instance, giving the rest of
the toolchain a consistent view of what was installed.

## Registry abstraction

Registries implement a shared interface with two capabilities: list packages
that are already available locally and enumerate remote versions that satisfy a
set of predicates. Implementations return lightweight package descriptors that
can later be resolved into local clones, alongside helpers for discovering the
modules contained within a package. Blush stores cached repositories under a
`git/<source>/<version>/` tree relative to the registry's root directory.

### Git registry provider

The built-in `GitRegistry` keeps its workspace inside the filesystem passed to
`PackageManager.New`. Local discovery scans each cached repository and returns
every version that has been cloned already. Remote discovery fetches tags from
the upstream Git repository, filters them by the provided predicates, sorts the
matches in descending version order, and offers them as resolvable packages.
Resolving a package clones the tag into a mangled `<source>/<version>/`
directory beneath the registry root.

Once a repository is available locally, the registry uses the filesystem module
discovery helper (`registry/fsmodule`) to list modules within the checkout. Each
module reports the `.blush` sources in its directory along with a logical URI
built from the package source and path.

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
