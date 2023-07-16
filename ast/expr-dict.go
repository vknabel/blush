package ast

var _ Expr = ExprDict{}

type ExprDict struct {
	Entries []ExprDictEntry
}

func MakeExprDict(entries []ExprDictEntry, source *Source) *ExprDict {
	return &ExprDict{
		Entries: entries,
	}
}

func (e ExprDict) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	for _, el := range e.Entries {
		el.Key.EnumerateNestedDecls(enumerate)
		el.Value.EnumerateNestedDecls(enumerate)
	}
}

type ExprDictEntry struct {
	Key   Expr
	Value Expr
}

func MakeExprDictEntry(key Expr, value Expr, source *Source) *ExprDictEntry {
	return &ExprDictEntry{
		Key:   key,
		Value: value,
	}
}
