// resolution is used to resolve the package path to the actual path
package registry

import (
	"context"
	"strings"

	"github.com/vknabel/blush/version"
)

// Provider is the registry for all packages in all versions.
// It is used to resolve the package path to the actual path.
//
// The expected folder structure in increasing priority is:
//
//	 $LITHIA_STDLIB/
//	 └── git/<package>/<version>/
//		 ├── Potfile
//	 	 └── <submodule>/
//	 $LITHIA_PACKAGES/
//	 └── git/<package>/<version>/
//		 ├── Potfile
//	 	 └── <submodule>/
//	 <package>
//	 ├── Potfile
//	 ├── <vendored-package>/
//	 │	 ├── Potfile
//	 │	 └── <submodule>/
//	 └── <submodule>/
//
// Each Potfile describes the package and its dependencies.
// The Potfile also declares the package name which are used for the imports.
// Each dependency can be renamed within the package.
type Provider interface {
	// Discover returns all packages in all versions that are available locally.
	Discover(ctx context.Context) ([]ResolvedPackage, error)
	// DiscoverPackageVersions returns packages with the given name constrained by the given predicates from remote.
	DiscoverPackageVersions(ctx context.Context, name string, preds ...version.Predicate) ([]Package, error)
}

type Package interface {
	Source() string
	Version() version.Version
	Resolve(ctx context.Context) (ResolvedPackage, error)
}

type ResolvedPackage interface {
	Package
	// ResolveModules discovers all nested modules.
	ResolveModules() ([]ResolvedModule, error)
}

type ResolvedModule interface {
	URI() LogicalURI
	Sources() ([]Source, error)
}

type Source interface {
	URI() LogicalURI
	Read() ([]byte, error)
}

type LogicalURI string

func (u LogicalURI) Join(segment string) LogicalURI {
	joined := u
	if !strings.HasSuffix(string(u), "/") {
		joined += "/"
	}
	return joined + LogicalURI(segment)
}
