package parser

import (
	"fmt"

	"github.com/vknabel/lithia/token"
)

type ParseError interface {
	error
	Source() *token.Source
	Summary() string
	Details() string
}

type UnexpectedTokenErr struct {
	Got      token.Token
	Expected []token.TokenType
}

// Error implements ParseError.
func (e UnexpectedTokenErr) Error() string {
	return fmt.Sprintf("unexpected %s %q, expected %s", e.Got.Type, e.Got.Literal, e.Expected)
}

func (e UnexpectedTokenErr) Summary() string {
	return fmt.Sprintf("syntax error")
}

func (e UnexpectedTokenErr) Details() string {
	return fmt.Sprintf("unexpected %s %q, expected %s", e.Got.Type, e.Got.Literal, e.Expected)
}

func (e UnexpectedTokenErr) Source() *token.Source {
	return e.Got.Source
}

func UnexpectedGot(got token.Token, expected ...token.TokenType) UnexpectedTokenErr {
	return UnexpectedTokenErr{got, expected}
}
