package runtime

import (
	"fmt"

	"github.com/vknabel/blush/ast"
)

var _ CallableRuntimeValue = &DataType{}

type DataType struct {
	symbol       *ast.Symbol
	FieldSymbols []*ast.Symbol
}

func MakeDataType(module *ast.SymbolTable, symbol *ast.Symbol) (*DataType, error) {
	decl, ok := symbol.Decl.(*ast.DeclData)
	if !ok {
		return nil, fmt.Errorf("declaration is not a DeclData, got %T", symbol.Decl)
	}
	fieldSymbols := make([]*ast.Symbol, len(decl.Fields))
	for i, f := range decl.Fields {
		for _, fsym := range symbol.ChildTable.Symbols {
			if fsym.Decl.DeclName().String() == f.DeclName().String() {
				fieldSymbols[i] = fsym
			}
		}
		return nil, fmt.Errorf("no symbol for field: %q", f.DeclName().String())
	}

	return &DataType{
		symbol:       symbol,
		FieldSymbols: fieldSymbols,
	}, nil
}

// Arity implements Callable.
func (dt *DataType) Arity() int {
	return len(dt.FieldSymbols)
}

// Call implements Callable.
func (dt *DataType) Call(args []RuntimeValue) RuntimeValue {
	val := &DataValue{
		TypeId: TypeId(dt.symbol.ConstantId),
		Fields: make(map[string]RuntimeValue, len(dt.FieldSymbols)),
	}
	for i, f := range dt.FieldSymbols {
		val.Fields[f.Decl.DeclName().String()] = args[i]
	}
	return val
}

// Inspect implements Callable.
func (dt *DataType) Inspect() string {
	return fmt.Sprintf("data %s", dt.symbol.Decl.DeclName())
}

// Lookup implements Callable.
func (dt *DataType) Lookup(name string) RuntimeValue {
	panic("unimplemented")
}

// TypeConstantId implements Callable.
func (dt *DataType) TypeConstantId() TypeId {
	return TypeId(dt.symbol.TypeSymbol.ConstantId)
}
