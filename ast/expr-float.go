package ast

var _ Expr = ExprFloat{}

type ExprFloat struct {
	Literal float64
}

func MakeExprFloat(literal float64, source *Source) *ExprFloat {
	return &ExprFloat{
		Literal: literal,
	}
}

func (e ExprFloat) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}
