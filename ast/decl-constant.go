package ast

import (
	"fmt"

	"github.com/vknabel/lithia/token"
)

var _ Decl = DeclVariable{}
var _ Overviewable = DeclVariable{}

type DeclVariable struct {
	Name  Identifier
	Value Expr
	Token token.Token

	Docs *Docs
}

// TokenLiteral implements Node
func (d DeclVariable) TokenLiteral() token.Token {
	return d.Token
}

// statementNode implements Statement
func (DeclVariable) statementNode() {}

// declarationNode implements Declaration
func (DeclVariable) declarationNode() {}

func (e DeclVariable) DeclName() Identifier {
	return e.Name
}

func (e DeclVariable) DeclOverview() string {
	return fmt.Sprintf("let %s", e.Name)
}

func (e DeclVariable) IsExportedDecl() bool {
	return true
}

func MakeDeclVariable(tok token.Token, name Identifier, value Expr) *DeclVariable {
	return &DeclVariable{
		Token: tok,
		Name:  name,
		Value: value,
	}
}

func (e DeclVariable) ProvidedDocs() *Docs {
	return e.Docs
}

func (e DeclVariable) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Value.EnumerateNestedDecls(enumerate)
}
