package ast

var _ Expr = ExprIdentifier{}

type ExprIdentifier struct {
	Name Identifier
}

func MakeExprIdentifier(name Identifier, source *Source) *ExprIdentifier {
	return &ExprIdentifier{
		Name: name,
	}
}

func (e ExprIdentifier) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}
