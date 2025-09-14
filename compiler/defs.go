package compiler

import (
	"github.com/vknabel/blush/op"
	"github.com/vknabel/blush/runtime"
)

type Compiler struct {
	instructions op.Instructions
	constants    []runtime.RuntimeValue
	plugins      *runtime.ExternPluginRegistry
}

func New() *Compiler {
	return &Compiler{
		instructions: op.Instructions{},
		constants:    []runtime.RuntimeValue{},
		plugins:      &runtime.ExternPluginRegistry{},
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

func (c *Compiler) emit(opcode op.Opcode, operands ...int) int {
	ins := op.Make(opcode, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	c.instructions = append(c.instructions, ins...)
	return len(c.instructions) - len(ins)
}

func (c *Compiler) addConstant(v runtime.RuntimeValue) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}
