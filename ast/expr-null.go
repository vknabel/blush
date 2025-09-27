package ast

import (
	"github.com/vknabel/blush/token"
)

var _ Expr = ExprNull{}

type ExprNull struct {
	Token token.Token
}

func MakeExprNull(token token.Token) *ExprNull {
	return &ExprNull{
		Token: token,
	}
}

// EnumerateChildNodes implements Expr.
func (ExprNull) EnumerateChildNodes(func(child Node)) {
	// No child nodes.
}

// TokenLiteral implements Expr.
func (n ExprNull) TokenLiteral() token.Token {
	return n.Token
}

// Expression implements Expr.
func (e ExprNull) Expression() string {
	return "null"
}
