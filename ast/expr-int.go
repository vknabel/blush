package ast

var _ Expr = ExprInt{}

type ExprInt struct {
	Literal int64
}

func MakeExprInt(literal int64, source *Source) *ExprInt {
	return &ExprInt{
		Literal: literal,
	}
}

func (e ExprInt) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}
