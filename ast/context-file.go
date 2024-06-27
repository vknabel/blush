package ast

import (
	"fmt"

	"github.com/vknabel/lithia/token"
)

var _ Node = &SourceFile{}

type SourceFile struct {
	Token      token.Token
	Path       string
	Statements []Statement
	Symbols    *SymbolTable
}

func MakeSourceFile(parent *SymbolTable, path string, token token.Token) *SourceFile {
	sf := &SourceFile{
		Token:      token,
		Path:       path,
		Statements: make([]Statement, 0),
	}
	sf.Symbols = MakeSymbolTable(parent, sf)
	sf.Symbols.delegatesExports = true
	return sf
}

func (sf *SourceFile) Add(globalStmt Statement) {
	if decl, ok := globalStmt.(Decl); ok {
		if sym := sf.Symbols.resolve(decl.DeclName().Value); sym == nil {
			sf.Symbols.Insert(decl)
		}
		return
	}
	sf.Statements = append(sf.Statements, globalStmt)
}

func (sf SourceFile) EnumerateChildNodes(action func(child Node)) {
	for _, sym := range sf.Symbols.Symbols {
		action(sym.Decl)
		sym.Decl.EnumerateChildNodes(action)
	}
	for _, sym := range sf.Symbols.Parent.Symbols {
		if sym.Decl == nil {
			fmt.Printf("SKIPPING UNDEFINED SYMBOL %+v", sym)
			continue
		}
		fmt.Printf("SYMBOL: %+v\nDECL: %T\n\n", sym, sym.Decl)
		action(sym.Decl)
		sym.Decl.EnumerateChildNodes(action)
	}

	for _, node := range sf.Statements {
		action(node)
		node.EnumerateChildNodes(action)
	}
}

// TokenLiteral implements Node.
func (sf *SourceFile) TokenLiteral() token.Token {
	return sf.Token
}
