package ast

import "github.com/vknabel/lithia/token"

type Node interface {
	TokenLiteral() token.Token
	// String() string
}

// Expressions produce values.
type Expression interface {
	Node
	expressionNode()
}

// Statements can be evaluated.
type Statement interface {
	Node
	statementNode()
}

// Delarations are statements that provide new bindings.
type Declaration interface {
	Node
	declarationNode()
}

type StatementDeclaration interface {
	Statement
	Declaration
}
