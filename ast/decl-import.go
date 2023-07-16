package ast

import (
	"fmt"
	"strings"

	"github.com/vknabel/lithia/token"
)

var _ Decl = DeclImport{}
var _ Overviewable = DeclImport{}

type DeclImport struct {
	Token      token.Token
	Alias      Identifier
	ModuleName ModuleName
	Members    []DeclImportMember
}

// TokenLiteral implements Node
func (d DeclImport) TokenLiteral() token.Token {
	return d.Token
}

// statementNode implements Statement
func (d DeclImport) statementNode() {}

// declarationNode implements Declaration
func (d DeclImport) declarationNode() {}

func (e DeclImport) DeclName() Identifier {
	return e.Alias
}

func (e DeclImport) DeclOverview() string {
	if e.Alias.Value != "" {
		return fmt.Sprintf("import %s = %s", e.Alias.Value, e.ModuleName)
	} else {
		return fmt.Sprintf("import %s", e.ModuleName)
	}
}

func (e DeclImport) IsExportedDecl() bool {
	return false
}

func (e *DeclImport) AddMember(member DeclImportMember) {
	e.Members = append(e.Members, member)
}

func MakeDeclImport(tok token.Token, segments []Identifier) *DeclImport {
	if len(segments) == 0 {
		panic("TODO: import declarations must consist of at least one identifier")
	}

	var nameBuilder strings.Builder
	for _, seg := range segments {
		nameBuilder.WriteString(seg.Value)
		nameBuilder.WriteRune('.')
	}
	moduleName := ModuleName(nameBuilder.String())
	moduleName = moduleName[:len(moduleName)-1]

	alias := Identifier(segments[len(segments)-1])
	return &DeclImport{
		Token:      tok,
		Alias:      alias,
		ModuleName: moduleName,
		Members:    make([]DeclImportMember, 0),
	}
}

func MakeDeclAliasImport(tok token.Token, alias Identifier, name ModuleName, source *Source) *DeclImport {
	return &DeclImport{
		Token:      tok,
		Alias:      alias,
		ModuleName: name,
		Members:    make([]DeclImportMember, 0),
	}
}

func (DeclImport) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}
