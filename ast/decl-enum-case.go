package ast

import "github.com/vknabel/lithia/token"

var _ Decl = DeclEnumCase{}

type DeclEnumCase struct {
	Token token.Token
	Name  Identifier

	Docs *Docs
}

// TokenLiteral implements Node
func (d DeclEnumCase) TokenLiteral() token.Token {
	return d.Token
}

// declarationNode implements Declaration
func (d DeclEnumCase) declarationNode() {}

func (e DeclEnumCase) DeclName() Identifier {
	return e.Name
}

func (e DeclEnumCase) IsExportedDecl() bool {
	return true
}

func MakeDeclEnumCase(tok token.Token, name Identifier) *DeclEnumCase {
	return &DeclEnumCase{
		Token: tok,
		Name:  name,
	}
}

func (decl DeclEnumCase) ProvidedDocs() *Docs {
	return decl.Docs
}

func (n DeclEnumCase) EnumerateChildNodes(action func(child Node)) {
	action(n.Name)
}
