package gitreg

import (
	"context"

	"github.com/vknabel/lithia/registry"
	"github.com/vknabel/lithia/version"
)

type localGitPackage struct {
	*remoteGitPackage
	localPath string
}

// Source implements registry.Package
func (p *localGitPackage) Source() string {
	return p.source
}

// Version implements registry.Package
func (p *localGitPackage) Version() version.Version {
	return p.version
}

// Resolve implements registry.Package
func (p *localGitPackage) Resolve(ctx context.Context) (registry.LocalPackage, error) {
	return p, nil
}

// LocalPath implements registry.LocalPackage
func (p *localGitPackage) LocalPath() string {
	return p.localPath
}
