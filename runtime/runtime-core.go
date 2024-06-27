package runtime

type TypeId uint32

type RuntimeValue interface {
	TypeConstantId() TypeId
	// for printing
	Inspect() string
	// Member access. If nil returned, this produces a runtime error.
	Lookup(name string) RuntimeValue
}

type CallableRuntimeValue interface {
	RuntimeValue
	Arity() int
	// Call expectes `len(args)` to exactly equal `Arity()`
	Call(args []RuntimeValue) RuntimeValue
}
