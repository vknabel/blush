package ast

import (
	"github.com/vknabel/lithia/token"
)

type AnnotationChain []*DeclAnnotationInstance

func MakeAnnotationChain(instances ...*DeclAnnotationInstance) AnnotationChain {
	return instances
}

// TokenLiteral implements Node
func (n AnnotationChain) TokenLiteral() token.Token {
	return n[0].TokenLiteral()
}

func (n AnnotationChain) EnumerateChildNodes(action func(child Node)) {
	for _, c := range n {
		action(c)
	}
}
