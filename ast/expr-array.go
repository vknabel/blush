package ast

import "github.com/vknabel/lithia/token"

var _ Expr = ExprArray{}

type ExprArray struct {
	Token    token.Token
	Elements []Expr
}

// TokenLiteral implements Expr.
func (e ExprArray) TokenLiteral() token.Token {
	return e.Token
}

func MakeExprArray(elements []Expr, token token.Token) *ExprArray {
	return &ExprArray{
		Elements: elements,
		Token:    token,
	}
}

func (e ExprArray) EnumerateChildNodes(enumerate func(Node)) {
	for _, el := range e.Elements {
		enumerate(el)
	}
}
