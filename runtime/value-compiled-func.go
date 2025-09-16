package runtime

import "github.com/vknabel/blush/op"

// CompiledFunction represents a blush function compiled to bytecode.
type CompiledFunction struct {
	Instructions  op.Instructions
	Constants     []RuntimeValue
	NumParameters int
}

// Ensure CompiledFunction implements RuntimeValue.
var _ RuntimeValue = (*CompiledFunction)(nil)

// Inspect implements RuntimeValue.
func (cf *CompiledFunction) Inspect() string {
	return "func"
}

// Lookup implements RuntimeValue.
func (cf *CompiledFunction) Lookup(name string) RuntimeValue {
	return nil
}

// TypeConstantId implements RuntimeValue.
func (cf *CompiledFunction) TypeConstantId() TypeId {
	return 0
}
