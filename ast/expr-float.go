package ast

import "github.com/vknabel/lithia/token"

var _ Expr = ExprFloat{}

type ExprFloat struct {
	Token   token.Token
	Literal float64
}

func MakeExprFloat(literal float64, token token.Token) *ExprFloat {
	return &ExprFloat{
		Literal: literal,
	}
}

// TokenLiteral implements Expr.
func (e ExprFloat) TokenLiteral() token.Token {
	return e.Token
}

func (e ExprFloat) EnumerateChildNodes(enumerate func(Node)) {
	// No child nodes.
}
