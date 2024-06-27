package ast

type Decl interface {
	Declaration
	DeclName() Identifier
	IsExportedDecl() bool
}

type DeclWithSymbols interface {
	Decl
	Symbols() *SymbolTable
}
