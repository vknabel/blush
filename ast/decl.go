package ast

type Decl interface {
	Declaration
	DeclName() Identifier
	IsExportedDecl() bool
}
