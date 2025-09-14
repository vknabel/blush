package pkgmanager

import (
	"github.com/go-git/go-billy/v5"
	"github.com/vknabel/blush/potfile"
	"github.com/vknabel/blush/registry"
	"github.com/vknabel/blush/registry/gitreg"
)

type PackageManager struct {
	registries []registry.Provider
}

func New(fs billy.Filesystem) (*PackageManager, error) {
	gitregfs, err := fs.Chroot("git")
	if err != nil {
		return nil, err
	}
	return &PackageManager{
		registries: []registry.Provider{
			gitreg.New(gitregfs),
		},
	}, nil
}

func (pm *PackageManager) Install(pot potfile.Potfile) *InstallationTask {
	return &InstallationTask{
		pkgmanager: pm,
		pot:        pot,
	}
}
