package ast

import (
	"github.com/vknabel/lithia/token"
)

var _ Expr = ExprFunc{}

type ExprFunc struct {
	Token      token.Token
	Name       string
	Parameters []DeclParameter
	Impl       Block
}

func MakeExprFunc(token token.Token, name string, parameters []DeclParameter, impl Block) *ExprFunc {
	return &ExprFunc{
		Token:      token,
		Name:       name,
		Parameters: parameters,
		Impl:       impl,
	}
}

// EnumerateChildNodes implements Expr.
func (n ExprFunc) EnumerateChildNodes(action func(child Node)) {
	for _, node := range n.Parameters {
		action(node)
		node.EnumerateChildNodes(action)
	}
	for _, node := range n.Impl {
		action(node)
		node.EnumerateChildNodes(action)
	}
}

// TokenLiteral implements Expr.
func (e ExprFunc) TokenLiteral() token.Token {
	return e.Token
}
