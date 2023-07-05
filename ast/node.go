package ast

import "github.com/vknabel/lithia/token"

type Node interface {
	Token() token.Token
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

type Declaration interface {
	Node
	declarationNode()
}
