package ast

var _ Expr = ExprTypeSwitch{}

type ExprTypeSwitch struct {
	Type      Expr
	CaseOrder []Identifier
	Cases     map[string]Expr
}

func MakeExprTypeSwitch(type_ Expr, source *Source) *ExprTypeSwitch {
	return &ExprTypeSwitch{
		Type:      type_,
		CaseOrder: make([]Identifier, 0),
		Cases:     make(map[string]Expr),
	}
}

func (e *ExprTypeSwitch) AddCase(key Identifier, value Expr) {
	e.CaseOrder = append(e.CaseOrder, key)
	e.Cases[key.Value] = value
}

func (e ExprTypeSwitch) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Type.EnumerateNestedDecls(enumerate)
	for _, ident := range e.CaseOrder {
		e.Cases[ident.Value].EnumerateNestedDecls(enumerate)
	}
}
