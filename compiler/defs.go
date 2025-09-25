package compiler

import (
	"github.com/vknabel/blush/ast"
	"github.com/vknabel/blush/op"
	"github.com/vknabel/blush/runtime"
)

type emittedInstruction struct {
	Opcode   op.Opcode
	Position int
}

type CompilationScope struct {
	instructions op.Instructions
	symbols      *ast.SymbolTable

	lastInstruction     emittedInstruction
	previousInstruction emittedInstruction
}

type Bytecode struct {
	Instructions op.Instructions
	Constants    []runtime.RuntimeValue
}

type Compiler struct {
	constants []runtime.RuntimeValue
	plugins   *runtime.ExternPluginRegistry

	scopes   []*CompilationScope
	scopeIdx int
}

func New() *Compiler {
	mainScope := &CompilationScope{
		instructions: op.Instructions{},
		symbols:      ast.MakeSymbolTable(nil, nil),
	}
	return &Compiler{
		constants: []runtime.RuntimeValue{},
		plugins:   &runtime.ExternPluginRegistry{},
		scopes:    []*CompilationScope{mainScope},
		scopeIdx:  0,
	}
}

func (c *Compiler) currentInstructions() op.Instructions {
	return c.scopes[c.scopeIdx].instructions
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) emit(opcode op.Opcode, operands ...int) int {
	ins := op.Make(opcode, operands...)
	pos := c.addInstruction(ins)

	c.scopes[c.scopeIdx].previousInstruction = c.scopes[c.scopeIdx].lastInstruction
	c.scopes[c.scopeIdx].lastInstruction = emittedInstruction{
		Opcode:   opcode,
		Position: pos,
	}
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	newPos := len(c.currentInstructions())
	c.scopes[c.scopeIdx].instructions = append(c.scopes[c.scopeIdx].instructions, ins...)
	return newPos
}

func (c *Compiler) addConstant(v runtime.RuntimeValue) int {
	c.constants = append(c.constants, v)
	// TODO: addConstant, what about types?
	return len(c.constants) - 1
}

func (c *Compiler) enterScope(syms *ast.SymbolTable) {
	c.scopes = append(c.scopes, &CompilationScope{
		instructions: op.Instructions{},
		symbols:      syms,
	})
	c.scopeIdx++
}

func (c *Compiler) leaveScope() *CompilationScope {
	scope := c.scopes[c.scopeIdx]
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIdx--
	return scope
}

func (c *Compiler) isLastInstruction(opcodes ...op.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	ins := c.scopes[c.scopeIdx].lastInstruction
	for _, opcode := range opcodes {
		if ins.Opcode == opcode {
			return true
		}
	}

	return false
}

func (c *Compiler) removeLastInstruction() emittedInstruction {
	last := c.scopes[c.scopeIdx].lastInstruction
	previous := c.scopes[c.scopeIdx].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIdx].instructions = new
	c.scopes[c.scopeIdx].lastInstruction = previous

	return last
}
