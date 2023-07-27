package token

type DecorativeTokenType string

type DecorativeToken struct {
	Type    DecorativeTokenType
	Literal string
}

const (
	COMMENT    DecorativeTokenType = "COMMENT"
	WHITESPACE DecorativeTokenType = "WHITESPACE"
)
