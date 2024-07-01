package vm

import (
	"fmt"

	"github.com/vknabel/lithia/op"
	"github.com/vknabel/lithia/runtime"
)

func (vm *VM) Run() error {
	var (
		ip   int
		code op.Opcode
	)

	for ip = 0; ip < len(vm.instructions); ip++ {
		code = op.Opcode(vm.instructions[ip])
		switch code {
		case op.Pop:
			vm.pop()

		case op.Const:
			idx := op.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[idx])
			if err != nil {
				return err
			}

		case op.Jump:
			pos := int(op.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case op.JumpFalse:
			pos := int(op.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			cond := vm.pop()

			if cond == runtime.Bool(false) {
				ip = pos - 1
			}
		case op.JumpTrue:
			pos := int(op.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			cond := vm.pop()

			if cond != runtime.Bool(false) {
				ip = pos - 1
			}

		case op.Add, op.Sub, op.Mul, op.Div:
			vm.numericBinaryOperation(code)

		default:
			def, err := op.Lookup(byte(code))
			if err != nil {
				return fmt.Errorf("unhandled opcode: %w", err)
			}
			return fmt.Errorf("unknown opcode %q", def.Name)
		}
	}
	return nil
}

func (vm *VM) push(val runtime.RuntimeValue) error {
	if vm.sp >= stackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = val
	vm.sp++
	return nil
}

func (vm *VM) pop() runtime.RuntimeValue {
	v := vm.stack[vm.sp-1]
	vm.sp--
	return v
}

func (vm *VM) numericBinaryOperation(operator op.Opcode) error {
	switch rhs := vm.pop().(type) {
	case runtime.Int:
		switch lhs := vm.pop().(type) {
		case runtime.Int:
			vm.numericBinaryOperationInt(operator, lhs, rhs)
		case runtime.Float:
			vm.numericBinaryOperationFloat(operator, lhs, runtime.Float(rhs))
		default:
			return fmt.Errorf("unsupported %T", lhs)
		}
	case runtime.Float:
		switch lhs := vm.pop().(type) {
		case runtime.Int:
			vm.numericBinaryOperationFloat(operator, runtime.Float(lhs), rhs)
		case runtime.Float:
			vm.numericBinaryOperationFloat(operator, lhs, rhs)
		default:
			return fmt.Errorf("unsupported %T", lhs)
		}
	default:
		return fmt.Errorf("unsupported %T", rhs)
	}
	return nil
}

func (vm *VM) numericBinaryOperationInt(operator op.Opcode, lhs, rhs runtime.Int) {
	switch operator {
	case op.Add:
		vm.push(lhs + rhs)
	case op.Sub:
		vm.push(lhs - rhs)
	case op.Mul:
		vm.push(lhs * rhs)
	case op.Div:
		vm.push(lhs / rhs)
	}
}
func (vm *VM) numericBinaryOperationFloat(operator op.Opcode, lhs, rhs runtime.Float) {
	switch operator {
	case op.Add:
		vm.push(lhs + rhs)
	case op.Sub:
		vm.push(lhs - rhs)
	case op.Mul:
		vm.push(lhs * rhs)
	case op.Div:
		vm.push(lhs / rhs)
	}
}
