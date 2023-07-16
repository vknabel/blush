package ast

var _ Expr = ExprArray{}

type ExprArray struct {
	Elements []Expr
}

func MakeExprArray(elements []Expr, source *Source) *ExprArray {
	return &ExprArray{
		Elements: elements,
	}
}

func (e ExprArray) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	for _, el := range e.Elements {
		el.EnumerateNestedDecls(enumerate)
	}
}
