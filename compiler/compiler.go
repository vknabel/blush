package compiler

import (
	"fmt"
	"math"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/op"
	"github.com/vknabel/lithia/token"
)

const (
	// A temporary address that acts placeholder.
	// Should be replaced by the actual address once known.
	placeholderJumpAddress = math.MinInt
)

func (c *Compiler) changeOperand(pos int, operand int) {
	opcode := op.Opcode(c.instructions[pos])
	patched := op.Make(opcode, operand)
	c.replaceInstruction(pos, patched)
}

func (c *Compiler) replaceInstruction(pos int, patched []byte) {
	for i := 0; i < len(patched); i++ {
		c.instructions[pos+i] = patched[i]
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.SourceFile:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.StmtExpr:
		err := c.Compile(node.Expr)
		if err != nil {
			return err
		}
		c.emit(op.Pop)
		return nil

	case ast.StmtIf:
		var (
			jumpNext int
			jumpEnds []int = make([]int, 0, 1+len(node.ElseIf))
			endPos   int
		)
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		jumpNext = c.emit(op.JumpFalse, placeholderJumpAddress)

		err = c.compileBlock(node.IfBlock)
		if err != nil {
			return err
		}

		jumpEnds = append(jumpEnds, c.emit(op.Jump, placeholderJumpAddress))

		for _, elseIf := range node.ElseIf {
			c.changeOperand(jumpNext, len(c.instructions))

			err = c.Compile(elseIf.Condition)
			if err != nil {
				return err
			}
			jumpNext = c.emit(op.JumpFalse, placeholderJumpAddress)

			err = c.compileBlock(elseIf.Block)
			if err != nil {
				return err
			}
			jumpEnds = append(jumpEnds, c.emit(op.Jump, placeholderJumpAddress))
		}
		c.changeOperand(jumpNext, len(c.instructions))

		err = c.compileBlock(node.ElseBlock)
		if err != nil {
			return err
		}

		endPos = len(c.instructions)
		for _, pos := range jumpEnds {
			c.changeOperand(pos, endPos)
		}
		return nil

	case ast.ExprIf:
		var (
			jumpNext int
			jumpEnds []int = make([]int, 0, 1+len(node.ElseIf))
			endPos   int
		)
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		jumpNext = c.emit(op.JumpFalse, placeholderJumpAddress)

		err = c.Compile(node.ThenExpr)
		if err != nil {
			return err
		}

		jumpEnds = append(jumpEnds, c.emit(op.Jump, placeholderJumpAddress))

		for _, elseIf := range node.ElseIf {
			c.changeOperand(jumpNext, len(c.instructions))

			err = c.Compile(elseIf.Condition)
			if err != nil {
				return err
			}
			jumpNext = c.emit(op.JumpFalse, placeholderJumpAddress)

			err = c.Compile(elseIf.Then)
			if err != nil {
				return err
			}
			jumpEnds = append(jumpEnds, c.emit(op.Jump, placeholderJumpAddress))
		}
		c.changeOperand(jumpNext, len(c.instructions))

		err = c.Compile(node.ElseExpr)
		if err != nil {
			return err
		}

		endPos = len(c.instructions)
		for _, pos := range jumpEnds {
			c.changeOperand(pos, endPos)
		}
		return nil

	case *ast.ExprOperatorBinary:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		// TODO: operators with lazy evaluation like && or ||

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator.Type {
		case token.PLUS:
			c.emit(op.Add)
			return nil
		case token.MINUS:
			c.emit(op.Sub)
			return nil
		case token.ASTERISK:
			c.emit(op.Mul)
			return nil
		case token.SLASH:
			c.emit(op.Div)
			return nil
		default:
			return fmt.Errorf("unknown infix operator %q", node.Operator.Literal)
		}

	case *ast.ExprInt:
		val := c.plugins.Prelude().Int(node.Literal)
		idx := c.addConstant(val)
		c.emit(op.Const, idx)
		return nil
	case *ast.ExprFloat:
		val := c.plugins.Prelude().Float(node.Literal)
		idx := c.addConstant(val)
		c.emit(op.Const, idx)
		return nil
	case *ast.ExprString:
		val := c.plugins.Prelude().String(node.Literal)
		idx := c.addConstant(val)
		c.emit(op.Const, idx)
		return nil

	default:
		return fmt.Errorf("unknown ast node %T", node)
	}
}

func (c *Compiler) compileBlock(block ast.Block) error {
	for _, stmt := range block {
		err := c.Compile(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
