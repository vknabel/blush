package potfile

import "github.com/vknabel/lithia/version"

type Potfile struct {
	Dependencies []Dependency
}

type Dependency struct {
	ImportName string
	Source     string
	Predicate  version.Predicate
}
