package vm

import (
	"github.com/vknabel/blush/compiler"
	"github.com/vknabel/blush/op"
	"github.com/vknabel/blush/runtime"
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

func newFrame(closure *runtime.Closure, basep int) *Frame {
	return &Frame{
		closure: closure,
		ip:      0,
		basep:   basep,
	}
}

func (f *Frame) Instructions() op.Instructions {
	return f.closure.Fn.Instructions
}

type VM struct {
	constants []runtime.RuntimeValue
	stack     []runtime.RuntimeValue
	sp        int
	frames    []*Frame
	framesIdx int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := runtime.MakeCompiledFunction(bytecode.Instructions, 0, nil)
	mainClosure := runtime.MakeClosure(mainFn, nil)

	frames := make([]*Frame, maxFrames)
	frames[0] = newFrame(mainClosure, 0)

	return &VM{
		stack:     make([]runtime.RuntimeValue, stackSize),
		sp:        0,
		constants: bytecode.Constants,
		frames:    frames,
		framesIdx: 1,
	}
}
func (vm *VM) LastPoppedStackElem() runtime.RuntimeValue {
	return vm.stack[vm.sp]
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIdx-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIdx] = f
	vm.framesIdx++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIdx--
	return vm.frames[vm.framesIdx]
}
