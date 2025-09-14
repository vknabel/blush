package runtime

import (
	"fmt"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/op"
)

var _ CallableRuntimeValue = CompiledFunction{}

type CompiledFunction struct {
	Intructions op.Instructions
	Locals      int
	Params      int
	Symbol      *ast.Symbol
}

func MakeCompiledFunction(
	intructions op.Instructions,
	locals int,
	params int,
	symbol *ast.Symbol,
) CompiledFunction {
	return CompiledFunction{
		Intructions: intructions,
		Locals:      locals,
		Params:      params,
		Symbol:      symbol,
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
