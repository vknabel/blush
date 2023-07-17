package parser

import (
	"fmt"

	"github.com/vknabel/lithia/token"
)

type ParseError interface {
	error
}

type UnexpectedTokenErr struct {
	Got      token.Token
	Expected []token.TokenType
}

// Error implements ParseError.
func (e UnexpectedTokenErr) Error() string {
	return fmt.Sprintf("unexpected %s %q, expected %s", e.Got.Type, e.Got.Literal, e.Expected)
}

func UnexpectedGot(got token.Token, expected ...token.TokenType) UnexpectedTokenErr {
	return UnexpectedTokenErr{got, expected}
}
