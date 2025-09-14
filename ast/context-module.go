package ast

import (
	"github.com/vknabel/blush/registry"
	"github.com/vknabel/blush/token"
)

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

func (m *ContextModule) TokenLiteral() token.Token {
	return token.Token{
		Type:    token.MODULE_DIRECTORY,
		Literal: string(m.Name),
		Source: &token.Source{
			File:   string(m.Name),
			Offset: 0,
		},
		Leading: []token.DecorativeToken{},
	}
}

func (m *ContextModule) EnumerateChildNodes(action func(child Node)) {
	for _, src := range m.Files {
		action(src)
	}
}
