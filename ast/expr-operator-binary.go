package ast

var _ Expr = ExprOperatorBinary{}

type ExprOperatorBinary struct {
	Operator OperatorBinary
	Left     Expr
	Right    Expr
}

func MakeExprOperatorBinary(operator OperatorBinary, left, right Expr, source *Source) *ExprOperatorBinary {
	return &ExprOperatorBinary{
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func (e ExprOperatorBinary) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Left.EnumerateNestedDecls(enumerate)
	e.Right.EnumerateNestedDecls(enumerate)
}
