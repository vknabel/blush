package runtime

import (
	"github.com/vknabel/lithia/ast"
)

// Declares some TypeId constants for the prelude data types.
// These are not guaranteed to be constant over versions and are not safe to serialize.
// They only offer fast creation for literals without excessive lookups.
// May change in the future.
const (
	typeIdArray TypeId = iota
	typeIdBool
	typeIdChar
	typeIdDict
	typeIdFloat
	typeIdFunc
	typeIdInt
	typeIdModule
	typeIdString
)

var _ ExternPlugin = &Prelude{}

type Prelude struct {
	symbols map[string]*ast.Symbol
}

// Bind implements runtime.ExternPlugin.
func (*Prelude) Bind(module *ast.SymbolTable, decl *ast.Symbol) RuntimeValue {
	switch decl.Name {
	case "Array",
		"Bool",
		"Char",
		"Dict",
		"Float",
		"Func",
		"Int",
		"Module",
		"String":
		return SimpleType{Decl: decl}
	case "Any":
		return MakeAnyType(decl)
	}
	return nil
}

func (p *Prelude) Bool(val bool) Bool                          { return Bool(val) }
func (p *Prelude) Array(val []RuntimeValue) Array              { return Array(val) }
func (p *Prelude) Char(val rune) Char                          { return Char(val) }
func (p *Prelude) Dict(val map[RuntimeValue]RuntimeValue) Dict { return Dict(val) }
func (p *Prelude) Float(val float64) Float                     { return Float(val) }
func (p *Prelude) Int(val int64) Int                           { return Int(val) }
func (p *Prelude) String(val string) String                    { return String(val) }
