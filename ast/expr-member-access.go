package ast

import "github.com/vknabel/lithia/token"

var _ Expr = ExprMemberAccess{}

type ExprMemberAccess struct {
	Target     Expr
	AccessPath []Identifier
}

func MakeExprMemberAccess(target Expr, accessPath []Identifier) *ExprMemberAccess {
	return &ExprMemberAccess{
		Target:     target,
		AccessPath: accessPath,
	}
}

// EnumerateChildNodes implements Expr.
func (n ExprMemberAccess) EnumerateChildNodes(action func(child Node)) {
	action(n.Target)
	for _, child := range n.AccessPath {
		action(child)
	}
}

// TokenLiteral implements Expr.
func (n ExprMemberAccess) TokenLiteral() token.Token {
	return n.Target.TokenLiteral()
}
