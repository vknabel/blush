package ast

var _ Expr = ExprMemberAccess{}

type ExprMemberAccess struct {
	Target     Expr
	AccessPath []Identifier
}

func MakeExprMemberAccess(target Expr, accessPath []Identifier, source *Source) *ExprMemberAccess {
	return &ExprMemberAccess{
		Target:     target,
		AccessPath: accessPath,
	}
}

func (e ExprMemberAccess) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Target.EnumerateNestedDecls(enumerate)
}
