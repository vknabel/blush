package runtime

import (
	"github.com/vknabel/blush/ast"
)

var _ RuntimeValue = SimpleType{}

type SimpleType struct {
	Decl *ast.Symbol
}

// Inspect implements runtime.RuntimeValue.
func (i SimpleType) Inspect() string {
	return "extern " + i.Decl.Name
}

// Lookup implements runtime.RuntimeValue.
func (i SimpleType) Lookup(name string) RuntimeValue {
	return nil
}

// TypeConstantId implements runtime.RuntimeValue.
func (i SimpleType) TypeConstantId() TypeId {
	return TypeId(i.Decl.ConstantId)
}
