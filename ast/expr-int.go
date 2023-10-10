package ast

import (
	"github.com/vknabel/lithia/token"
)

var _ Expr = ExprInt{}

type ExprInt struct {
	Literal int64
	Token   token.Token
}

func MakeExprInt(literal int64, token token.Token) *ExprInt {
	return &ExprInt{
		Literal: literal,
		Token:   token,
	}
}

// EnumerateChildNodes implements Expr.
func (ExprInt) EnumerateChildNodes(func(child Node)) {
	// No child nodes.
}

// TokenLiteral implements Expr.
func (n ExprInt) TokenLiteral() token.Token {
	return n.Token
}
