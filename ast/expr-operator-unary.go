package ast

import "github.com/vknabel/lithia/token"

var _ Expr = ExprOperatorUnary{}

type ExprOperatorUnary struct {
	Operator OperatorUnary
	Expr     Expr
}

func MakeExprOperatorUnary(operator OperatorUnary, expr Expr) *ExprOperatorUnary {
	return &ExprOperatorUnary{
		Operator: operator,
		Expr:     expr,
	}
}

// EnumerateChildNodes implements Expr.
func (n ExprOperatorUnary) EnumerateChildNodes(action func(child Node)) {
	action(n.Operator)
	action(n.Expr)
}

// TokenLiteral implements Expr.
func (n ExprOperatorUnary) TokenLiteral() token.Token {
	return n.Operator.TokenLiteral()
}
