package ast

import (
	"fmt"
	"github.com/vknabel/lithia/token"
	"strings"
)

var _ Decl = DeclFunc{}
var _ Overviewable = DeclFunc{}

type DeclFunc struct {
	Token token.Token
	Name  Identifier
	Impl  *ExprFunc

	Docs *Docs
}

// TokenLiteral implements Decl.
func (d DeclFunc) TokenLiteral() token.Token {
	return d.Token
}

// declarationNode implements Decl.
func (DeclFunc) declarationNode() {}

// statementNode implements Statement.
func (DeclFunc) statementNode() {}

func (e DeclFunc) DeclName() Identifier {
	return e.Name
}

func (e DeclFunc) DeclOverview() string {
	if len(e.Impl.Parameters) == 0 {
		return fmt.Sprintf("func %s { => }", e.Name)
	}
	paramNames := make([]string, len(e.Impl.Parameters))
	for i, param := range e.Impl.Parameters {
		paramNames[i] = string(param.Name.Value)
	}
	return fmt.Sprintf("func %s { %s => }", e.Name, strings.Join(paramNames, ", "))
}

func (e DeclFunc) IsExportedDecl() bool {
	return true
}

func MakeDeclFunc(tok token.Token, name Identifier, impl *ExprFunc, source *Source) *DeclFunc {
	return &DeclFunc{
		Token: tok,
		Name:  name,
		Impl:  impl,
	}
}

func (decl DeclFunc) ProvidedDocs() *Docs {
	return decl.Docs
}

func (decl DeclFunc) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	decl.Impl.EnumerateNestedDecls(enumerate)
}
