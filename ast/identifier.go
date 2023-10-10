package ast

import (
	"github.com/vknabel/lithia/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func MakeIdentifier(tok token.Token) Identifier {
	return Identifier{tok, tok.Literal}
}

func (Identifier) expressionNode() {}

func (n Identifier) String() string {
	return n.Value
}

func (n Identifier) TokenLiteral() token.Token {
	return n.Token
}

// EnumerateChildNodes implements Node.
func (n Identifier) EnumerateChildNodes(action func(child Node)) {}
