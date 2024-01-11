package ast

import (
	"fmt"
)

type SourceFile struct {
	Path         string
	Imports      []ModuleName
	Declarations []Decl
	Statements   []Expr
}

func MakeSourceFile(path string) *SourceFile {
	return &SourceFile{
		Path:         path,
		Imports:      make([]ModuleName, 0),
		Declarations: make([]Decl, 0),
		Statements:   make([]Expr, 0),
	}
}

func (sf *SourceFile) Add(globalStmt Statement) {
	switch node := globalStmt.(type) {
	case Decl:
		sf.Declarations = append(sf.Declarations, node)

	default:
		panic(fmt.Sprintf("TODO: unknown global statement %T", node))
	}
}

func (sf *SourceFile) AddDecl(decl Decl) {
	if decl == nil {
		return
	}
	if importDecl, ok := decl.(DeclImport); ok {
		sf.Imports = append(sf.Imports, importDecl.ModuleName)
	}
	sf.Declarations = append(sf.Declarations, decl)
}

func (sf *SourceFile) AddExpr(expr Expr) {
	if expr == nil {
		return
	}
	sf.Statements = append(sf.Statements, expr)
}

func (sf *SourceFile) ExportedDeclarations() []Decl {
	decls := make([]Decl, 0)
	for _, decl := range sf.Declarations {
		if decl.IsExportedDecl() {
			decls = append(decls, decl)
		}
	}
	return decls
}

func (sf SourceFile) EnumerateChildNodes(action func(child Node)) {
	for _, node := range sf.Declarations {
		action(node)
		node.EnumerateChildNodes(action)
	}
	for _, node := range sf.Statements {
		action(node)
		node.EnumerateChildNodes(action)
	}
}
