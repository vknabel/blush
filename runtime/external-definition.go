package runtime

import "github.com/vknabel/lithia/ast"

type ExternPlugin interface {
	Bind(module *ast.SymbolTable, decl *ast.Symbol) RuntimeValue
}
