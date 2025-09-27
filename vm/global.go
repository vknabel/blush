package vm

import (
	"errors"
	"sync/atomic"

	goruntime "runtime"

	"github.com/vknabel/blush/runtime"
)

type TaskId uint64

const (
	globalSlotStateUninitialized uint32 = iota
	globalSlotStateInitializing
	globalSlotStateInitialized
)

var errRecursiveVariableInitialization = errors.New("recursive initialization of global variable")

type Global struct {
	state uint32
	owner uint64
	init  func(TaskId) (runtime.RuntimeValue, error)
	value runtime.RuntimeValue
}

func MakeGlobal(init func(TaskId) (runtime.RuntimeValue, error)) *Global {
	return &Global{
		state: globalSlotStateUninitialized,
		init:  init,
	}
}

func (s *Global) Get(owner TaskId) (runtime.RuntimeValue, error) {
	var spin int

	for {
		state := atomic.LoadUint32(&s.state)
		switch state {
		case globalSlotStateInitialized:
			return s.value, nil

		case globalSlotStateUninitialized:
			if atomic.CompareAndSwapUint32(&s.state, globalSlotStateUninitialized, globalSlotStateInitializing) {
				atomic.StoreUint64(&s.owner, uint64(owner))
				v, err := s.init(owner)
				s.value = v
				s.init = nil
				atomic.StoreUint32(&s.state, globalSlotStateInitialized)
				return v, err
			}

		case globalSlotStateInitializing:
			spin++

			if spin < 100 {
				continue
			}

			if uint64(owner) == atomic.LoadUint64(&s.owner) {
				return nil, errRecursiveVariableInitialization
			}

			// Yield to avoid busy-waiting too long
			spin = 0
			goruntime.Gosched()
		}
	}
}

func (s *Global) Set(owner TaskId, v runtime.RuntimeValue) error {
	var spin int

	for {
		state := atomic.LoadUint32(&s.state)
		switch state {
		case globalSlotStateInitialized:
			s.value = v

		case globalSlotStateUninitialized:
			if atomic.CompareAndSwapUint32(&s.state, globalSlotStateUninitialized, globalSlotStateInitializing) {
				s.value = v
				s.init = nil
				atomic.StoreUint32(&s.state, globalSlotStateInitialized)
				return nil
			}

		case globalSlotStateInitializing:
			spin++

			if spin < 100 {
				continue
			}

			if uint64(owner) == atomic.LoadUint64(&s.owner) {
				return errRecursiveVariableInitialization
			}

			// Yield to avoid busy-waiting too long
			spin = 0
			goruntime.Gosched()
		}
	}
}
