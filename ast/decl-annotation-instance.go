package ast

import (
	"github.com/vknabel/lithia/token"
)

type AnnotationInstance struct {
	Token     token.Token
	Reference StaticReference
	Arguments []Expr
}

// TokenLiteral implements Node
func (n AnnotationInstance) TokenLiteral() token.Token {
	return n.Token
}

func MakeAnnotationInstance(tok token.Token, ref StaticReference) *AnnotationInstance {
	return &AnnotationInstance{tok, ref, nil}
}

func (n AnnotationInstance) AddArgument(arg Expr) {
	n.Arguments = append(n.Arguments, arg)
}

func (n AnnotationInstance) EnumerateChildNodes(action func(child Node)) {
	action(n.Reference)
	for _, argument := range n.Arguments {
		action(argument)
	}
}
