package ast

type Decl interface {
	Declaration
	DeclName() Identifier
	IsExportedDecl() bool

	EnumerateNestedDecls(enumerate func(interface{}, []Decl))
}
