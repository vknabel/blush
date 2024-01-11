package ast

import "github.com/vknabel/lithia/token"

var _ Statement = StmtReturn{}

type StmtReturn struct {
	Token token.Token
	Expr  Expr
}

func MakeStmtReturn(t token.Token, expr Expr) StmtReturn {
	return StmtReturn{
		Token: t,
		Expr:  expr,
	}
}

// EnumerateChildNodes implements Statement.
func (s StmtReturn) EnumerateChildNodes(action func(child Node)) {
	action(s.Expr)
}

// TokenLiteral implements Statement.
func (s StmtReturn) TokenLiteral() token.Token {
	return s.Token
}

// statementNode implements Statement.
func (s StmtReturn) statementNode() {}
