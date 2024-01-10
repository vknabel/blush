package ast

import (
	"github.com/vknabel/lithia/token"
)

type AnnotationChain struct {
	Token     token.Token
	Reference StaticReference
	Arguments []Expr
}

// TokenLiteral implements Node
func (n AnnotationChain) TokenLiteral() token.Token {
	return n.Token
}

func MakeAnnotationChain(tok token.Token, ref StaticReference) *AnnotationChain {
	return &AnnotationChain{tok, ref, nil}
}

func (n AnnotationChain) AddArgument(arg Expr) {
	n.Arguments = append(n.Arguments, arg)
}

func (n AnnotationChain) EnumerateChildNodes(action func(child Node)) {
	action(n.Reference)
	for _, argument := range n.Arguments {
		action(argument)
	}
}
