package parser

import (
	"fmt"

	"github.com/vknabel/lithia/token"
)

type ParseError interface {
	error
}

type UnexpectedToken struct {
	Expected token.TokenType
	Got      token.Token
}

// Error implements ParseError.
func (e UnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected %s %q, expected %s", e.Got.Type, e.Got.Literal, e.Expected)
}
