package ast

var _ Expr = ExprGroup{}

type ExprGroup struct {
	Expr Expr
}

func MakeExprGroup(expr Expr, source *Source) *ExprGroup {
	return &ExprGroup{
		Expr: expr,
	}
}

func (e ExprGroup) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Expr.EnumerateNestedDecls(enumerate)
}
