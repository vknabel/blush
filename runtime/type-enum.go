package runtime

import (
	"fmt"

	"github.com/vknabel/lithia/ast"
)

var _ RuntimeValue = &EnumType{}

type EnumType struct {
	symbol *ast.Symbol
}

// Inspect implements RuntimeValue.
func (et *EnumType) Inspect() string {
	return fmt.Sprintf("data %s", et.symbol.Decl.DeclName())
}

// Lookup implements RuntimeValue.
func (*EnumType) Lookup(name string) RuntimeValue {
	return nil
}

// TypeConstantId implements RuntimeValue.
func (et *EnumType) TypeConstantId() TypeId {
	return TypeId(et.symbol.TypeSymbol.ConstantId)
}
