package ast

var _ Expr = ExprOperatorUnary{}

type ExprOperatorUnary struct {
	Operator OperatorUnary
	Expr     Expr
}

func MakeExprOperatorUnary(operator OperatorUnary, expr Expr, source *Source) *ExprOperatorUnary {
	return &ExprOperatorUnary{
		Operator: operator,
		Expr:     expr,
	}
}

func (e ExprOperatorUnary) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Expr.EnumerateNestedDecls(enumerate)
}
