package vm

import (
	"github.com/vknabel/lithia/compiler"
	"github.com/vknabel/lithia/op"
	"github.com/vknabel/lithia/runtime"
)

const stackSize = 2048
const globalSize = 65536
const maxFrames = 1024

type VM struct {
	constants    []runtime.RuntimeValue
	stack        []runtime.RuntimeValue
	sp           int
	instructions op.Instructions
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		stack:        make([]runtime.RuntimeValue, stackSize),
		sp:           0,
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
	}
}
func (vm *VM) LastPoppedStackElem() runtime.RuntimeValue {
	return vm.stack[vm.sp]
}
