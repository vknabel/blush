package runtime

import "fmt"

type DataValue struct {
	TypeId TypeId
	Values []RuntimeValue
	Fields map[string]int
}

func MakeDataValue(dt *DataType, values []RuntimeValue) *DataValue {
	fields := make(map[string]int, len(dt.FieldSymbols))
	for i, f := range dt.FieldSymbols {
		fields[f.Name] = i
	}
	return &DataValue{
		TypeId: TypeId(*dt.Symbol.ConstantId),
		Fields: fields,
		Values: values,
	}
}

// Inspect implements RuntimeValue.
func (dv *DataValue) Inspect() string {
	return fmt.Sprintf("data #%d { %+v }", dv.TypeId, dv.Fields)
}

// Lookup implements RuntimeValue.
func (dv *DataValue) Lookup(name string) RuntimeValue {
	return dv.Values[dv.Fields[name]]
}

// TypeConstantId implements RuntimeValue.
func (dv *DataValue) TypeConstantId() TypeId {
	return dv.TypeId
}
