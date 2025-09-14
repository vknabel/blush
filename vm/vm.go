package vm

import (
	"github.com/vknabel/blush/compiler"
	"github.com/vknabel/blush/op"
	"github.com/vknabel/blush/runtime"
)

const stackSize = 2048
const globalSize = 65536
const maxFrames = 1024

type VM struct {
	constants    []runtime.RuntimeValue
	stack        []runtime.RuntimeValue
	sp           int
	globals      []runtime.RuntimeValue
	instructions op.Instructions
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		stack:        make([]runtime.RuntimeValue, stackSize),
		sp:           0,
		globals:      make([]runtime.RuntimeValue, globalSize),
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
	}
}
func (vm *VM) LastPoppedStackElem() runtime.RuntimeValue {
	return vm.stack[vm.sp]
}
