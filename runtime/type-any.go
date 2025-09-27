package runtime

import "github.com/vknabel/blush/ast"

var _ RuntimeValue = &AnyType{}

type AnyType struct {
	symbol *ast.Symbol
}

func MakeAnyType(symbol *ast.Symbol) *AnyType {
	return &AnyType{symbol}
}

// Inspect implements RuntimeValue.
func (*AnyType) Inspect() string {
	return "extern Any"
}

// Lookup implements RuntimeValue.
func (*AnyType) Lookup(name string) RuntimeValue {
	return nil
}

// TypeConstantId implements RuntimeValue.
func (at *AnyType) TypeConstantId() TypeId {
	return TypeId(*at.symbol.ConstantId)
}
