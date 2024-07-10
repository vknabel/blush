package runtime

var _ RuntimeValue = Void{}

type Void struct{}

// Inspect implements RuntimeValue.
func (v Void) Inspect() string {
	return "void"
}

// Lookup implements RuntimeValue.
func (v Void) Lookup(name string) RuntimeValue {
	panic("unimplemented")
}

// TypeConstantId implements RuntimeValue.
func (v Void) TypeConstantId() TypeId {
	return typeIdVoid
}
