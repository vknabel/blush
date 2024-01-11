package ast

import (
	"github.com/vknabel/lithia/token"
)

var _ Expr = ExprGroup{}

type ExprGroup struct {
	Token token.Token
	Expr  Expr
}

func MakeExprGroup(token token.Token, expr Expr) *ExprGroup {
	return &ExprGroup{
		Token: token,
		Expr:  expr,
	}
}

// EnumerateChildNodes implements Expr.
func (n ExprGroup) EnumerateChildNodes(action func(child Node)) {
	action(n.Expr)
}

// TokenLiteral implements Expr.
func (e ExprGroup) TokenLiteral() token.Token {
	return e.Token
}
