package ast

import (
	"github.com/vknabel/lithia/token"
)

var _ Expr = ExprInvocation{}

type ExprInvocation struct {
	Function  Expr
	Arguments []Expr
}

func MakeExprInvocation(function Expr, source *Source) *ExprInvocation {
	return &ExprInvocation{
		Function: function,
	}
}

func (e *ExprInvocation) AddArgument(argument Expr) {
	e.Arguments = append(e.Arguments, argument)
}

// EnumerateChildNodes implements Expr.
func (n ExprInvocation) EnumerateChildNodes(action func(child Node)) {
	action(n.Function)
	for _, argument := range n.Arguments {
		action(argument)
	}
}

// TokenLiteral implements Expr.
func (n ExprInvocation) TokenLiteral() token.Token {
	return n.Function.TokenLiteral()
}
