package runtime

var _ RuntimeValue = Null{}

type Null struct{}

// Inspect implements RuntimeValue.
func (n Null) Inspect() string {
	return "null"
}

// Lookup implements RuntimeValue.
func (n Null) Lookup(name string) RuntimeValue {
	return nil
}

// TypeConstantId implements RuntimeValue.
func (n Null) TypeConstantId() TypeId {
	return typeIdNull
}
