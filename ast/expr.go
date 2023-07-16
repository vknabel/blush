package ast

type Expr interface {
	EnumerateNestedDecls(enumerate func(interface{}, []Decl))
}
