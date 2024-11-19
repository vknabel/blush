package ast

import "github.com/vknabel/lithia/registry"

type ModuleName StaticReference

type ContextModule struct {
	Name    registry.LogicalURI
	Symbols *SymbolTable

	Files []*SourceFile
}

func MakeContextModule(name registry.LogicalURI) *ContextModule {
	m := &ContextModule{
		Name:  name,
		Files: []*SourceFile{},
	}
	m.Symbols = MakeSymbolTable(m.Symbols.Parent, nil)
	return m
}

func (m *ContextModule) AddSourceFile(sourceFile *SourceFile) {
	m.Files = append(m.Files, sourceFile)
}
