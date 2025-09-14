package vm

import (
	"github.com/vknabel/lithia/compiler"
	"github.com/vknabel/lithia/op"
	"github.com/vknabel/lithia/runtime"
)

const (
	stackSize  = 2048
	globalSize = 65536
	maxFrames  = 1024
)

type Frame struct {
	closure *runtime.Closure
	ip      int
	basep   int
}

type VM struct {
	constants    []runtime.RuntimeValue
	stack        []runtime.RuntimeValue
	frames       []Frame
	sp           int
	instructions op.Instructions
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		stack:        make([]runtime.RuntimeValue, stackSize),
		frames:       make([]Frame, maxFrames),
		sp:           0,
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
	}
}
func (vm *VM) LastPoppedStackElem() runtime.RuntimeValue {
	return vm.stack[vm.sp]
}
