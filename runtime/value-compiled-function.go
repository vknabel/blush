package runtime

import (
	"fmt"

	"github.com/vknabel/blush/ast"
	"github.com/vknabel/blush/op"
)

var _ CallableRuntimeValue = CompiledFunction{}

type CompiledFunction struct {
	Instructions op.Instructions
	Locals       int
	Params       int
	Symbol       *ast.Symbol
}

func MakeCompiledFunction(
	instructions op.Instructions,
	locals int,
	params int,
	symbol *ast.Symbol,
) CompiledFunction {
	return CompiledFunction{
		Instructions: instructions,
		Locals:       locals,
		Params:       params,
		Symbol:       symbol,
	}
}

// Arity implements CallableRuntimeValue.
func (c CompiledFunction) Arity() int {
	return c.Params
}

// Inspect implements CallableRuntimeValue.
func (c CompiledFunction) Inspect() string {
	return fmt.Sprintf("func %s(#%d)", c.Symbol.Decl.DeclName(), c.Arity())
}

// Lookup implements CallableRuntimeValue.
func (c CompiledFunction) Lookup(name string) RuntimeValue {
	if name == "arity" {
		return Int(c.Arity())
	}
	return nil
}

// TypeConstantId implements CallableRuntimeValue.
func (c CompiledFunction) TypeConstantId() TypeId {
	return TypeId(c.Symbol.TypeSymbol.ConstantId)
}
