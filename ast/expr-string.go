package ast

import "github.com/vknabel/lithia/token"

var _ Expr = ExprString{}

type ExprString struct {
	Token   token.Token
	Literal string
}

func MakeExprString(literal string, token token.Token) *ExprString {
	return &ExprString{
		Literal: literal,
		Token:   token,
	}
}

func (e ExprString) TokenLiteral() token.Token {
	return e.Token
}

func (e ExprString) EnumerateChildNodes(enumerate func(Node)) {
	// No child nodes.
}
