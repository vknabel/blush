package prelude

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.ExternPlugin = &prelude{}

type prelude struct{}

// Bind implements runtime.ExternPlugin.
func (*prelude) Bind(module *ast.SymbolTable, decl *ast.Symbol) runtime.RuntimeValue {
	panic("unimplemented")
}
