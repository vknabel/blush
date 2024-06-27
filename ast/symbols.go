package ast

import (
	"errors"
	"fmt"
)

var (
	errSymbolAlreadyDefinedInSameScope = errors.New("symbol already defined")
)

type Symbol struct {
	Name       string
	Decl       Decl
	Usages     []SymbolUsage
	ChildTable *SymbolTable
	Errs       []error

	// Filled by later phases

	ConstantId uint16
	TypeSymbol *Symbol
}

type SymbolUsage struct {
	Node             Node
	typeRequirements []SymbolRequirement
	Errs             []error
}

type SymbolRequirement interface{}
type RequireAnnotation *DeclAnnotationInstance
type RequireStaticRef struct {
	StaticReference
	ResolveRequirements SymbolRequirement
}

type SymbolTable struct {
	Parent   *SymbolTable
	OpenedBy Node
	Symbols  map[string]*Symbol

	functionCounter  int
	delegatesExports bool
}

func MakeSymbolTable(parent *SymbolTable, declaringNode Node) *SymbolTable {
	return &SymbolTable{
		Parent:   parent,
		OpenedBy: declaringNode,
		Symbols:  map[string]*Symbol{},
	}
}

func (st *SymbolTable) Name() string {
	var prefix string
	if st.Parent != nil {
		prefix = st.Parent.Name() + "->"
	}
	var name string
	switch n := st.OpenedBy.(type) {
	case Decl:
		name = n.DeclName().String()
	case ExprFunc:
		name = n.Name
	case Identifier:
		name = n.Value
	case *SourceFile:
		name = n.Path
	default:
		name = fmt.Sprintf("%T", st.OpenedBy)
	}
	return prefix + name
}

func (st *SymbolTable) Insert(decl Decl) *Symbol {
	if st.delegatesExports && decl.IsExportedDecl() {
		return st.Parent.Insert(decl)
	}
	name := decl.DeclName().Value
	if sym, ok := st.Symbols[name]; ok {
		sym.Errs = append(sym.Errs, errSymbolAlreadyDefinedInSameScope)
		sym.Usages = append(sym.Usages, SymbolUsage{
			Node:             decl,
			typeRequirements: nil,
			Errs:             []error{errSymbolAlreadyDefinedInSameScope},
		})
		return sym
	}
	sym := &Symbol{
		Name: decl.DeclName().Value,
		Decl: decl,
	}
	st.Symbols[name] = sym
	return sym
}

func (st *SymbolTable) addSymbol(symbol Symbol) *Symbol {
	if st.delegatesExports == true && (symbol.Decl == nil || symbol.Decl.IsExportedDecl()) {
		return st.Parent.addSymbol(symbol)
	}
	ref := &symbol
	st.Symbols[symbol.Name] = ref
	return ref
}

func (st *SymbolTable) resolve(name string) *Symbol {
	if st == nil {
		return nil
	}
	if sym, ok := st.Symbols[name]; ok {
		return sym
	}
	if sym := st.Parent.resolve(name); sym != nil {
		return sym
	}
	return nil
}

func (st *SymbolTable) Lookup(name string, fromNode Node, requirements ...SymbolRequirement) *Symbol {
	usage := SymbolUsage{
		Node:             fromNode,
		typeRequirements: requirements,
	}
	if sym := st.resolve(name); sym != nil {
		sym.Usages = append(sym.Usages, usage)
		return sym
	}
	return st.addSymbol(
		Symbol{
			Name:   name,
			Decl:   nil,
			Usages: []SymbolUsage{usage},
		},
	)
}

func (st *SymbolTable) LookupIdentifier(name Identifier, requirements ...SymbolRequirement) *Symbol {
	usage := SymbolUsage{
		Node:             name,
		typeRequirements: requirements,
	}
	if sym := st.resolve(name.Value); sym != nil {
		sym.Usages = append(sym.Usages, usage)
		return sym
	}
	return st.addSymbol(
		Symbol{
			Name:   name.Value,
			Decl:   nil,
			Usages: []SymbolUsage{usage},
		},
	)
}

func (st *SymbolTable) LookupRef(ref StaticReference, requirements ...SymbolRequirement) *Symbol {
	if len(ref) == 1 {
		return st.LookupIdentifier(ref[0], requirements...)
	}

	name := ref[0]
	usage := SymbolUsage{
		Node:             name,
		typeRequirements: append(requirements, RequireStaticRef{ref[1:], requirements}),
	}
	if sym := st.resolve(name.Value); sym != nil {
		if sym.ChildTable != nil {
			return sym.ChildTable.LookupRef(ref[1:])
		}
		if sym.Decl != nil {
			usage.Errs = append(usage.Errs, fmt.Errorf("expected to have member %s", ref[1:]))
		}
		sym.Usages = append(sym.Usages, usage)
		return sym
	}
	return st.addSymbol(
		Symbol{
			Name:   name.Value,
			Decl:   nil,
			Usages: []SymbolUsage{usage},
		},
	)
}

func (st *SymbolTable) NextAnonymousFunctionName() string {
	st.functionCounter++
	return fmt.Sprintf("func#%d", st.functionCounter)
}
