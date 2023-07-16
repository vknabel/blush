package ast

var _ Expr = ExprString{}

type ExprString struct {
	Literal string
}

func MakeExprString(literal string, source *Source) *ExprString {
	return &ExprString{
		Literal: literal,
	}
}

func (e ExprString) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}
