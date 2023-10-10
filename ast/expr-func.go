package ast

import (
	"github.com/vknabel/lithia/token"
)

var _ Expr = ExprFunc{}

type ExprFunc struct {
	Token        token.Token
	Name         string
	Parameters   []DeclParameter
	Declarations []Decl
	Expressions  []Expr
}

func MakeExprFunc(name string, parameters []DeclParameter, token token.Token) *ExprFunc {
	return &ExprFunc{
		Token:        token,
		Name:         name,
		Parameters:   parameters,
		Declarations: []Decl{},
		Expressions:  []Expr{},
	}
}

func (e *ExprFunc) AddDecl(decl Decl) {
	e.Declarations = append(e.Declarations, decl)
}

func (e *ExprFunc) AddExpr(expr Expr) {
	e.Expressions = append(e.Expressions, expr)
}

// EnumerateChildNodes implements Expr.
func (n ExprFunc) EnumerateChildNodes(action func(child Node)) {
	for _, node := range n.Parameters {
		action(node)
	}
	for _, node := range n.Declarations {
		action(node)
	}
	for _, node := range n.Expressions {
		action(node)
	}
}

// TokenLiteral implements Expr.
func (e ExprFunc) TokenLiteral() token.Token {
	return e.Token
}
