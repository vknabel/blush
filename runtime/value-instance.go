package runtime

import "fmt"

type DataValue struct {
	TypeId TypeId
	Fields map[string]RuntimeValue
}

// Inspect implements RuntimeValue.
func (dv *DataValue) Inspect() string {
	return fmt.Sprintf("data #%d { %+v }", dv.TypeId, dv.Fields)
}

// Lookup implements RuntimeValue.
func (dv *DataValue) Lookup(name string) RuntimeValue {
	panic("unimplemented")
}

// TypeConstantId implements RuntimeValue.
func (dv *DataValue) TypeConstantId() TypeId {
	return dv.TypeId
}
