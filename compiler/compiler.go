package compiler

import (
	"github.com/vknabel/lithia/op"
	"github.com/vknabel/lithia/runtime"
)

type Compiler struct {
	instructions op.Instructions
	constants    []runtime.RuntimeValue
}

func New() *Compiler {
	return &Compiler{
		instructions: op.Instructions{},
		constants:    []runtime.RuntimeValue{},
	}
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions op.Instructions
	Constants    []runtime.RuntimeValue
}
