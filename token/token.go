package token

type TokenType string

// TODO: add whitespace and comments to tokens
// This allows the parser to ignore all comments
// and the ast may also drop docs completely
type Token struct {
	Type    TokenType
	Literal string
	Source  *Source
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "COMMENT"

	// Identifiers + literals
	IDENT  = "IDENT"
	STRING = "STRING"
	INT    = "INT"
	FLOAT  = "FLOAT"

	// Operators
	BANG     = "!"
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	LT  = "<"
	GT  = ">"
	EQ  = "=="
	NEQ = "!="
	GTE = ">="
	LTE = "<="
	AND = "&&"
	OR  = "||"

	// Delimiters
	ASSIGN   = "="
	ARROW    = "=>"
	COLON    = ":"
	DOT      = "."
	COMMA    = ","
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	AT       = "@"

	// KEYWORDS
	MODULE   = "MODULE"
	IMPORT   = "IMPORT"
	ENUM     = "ENUM"
	DATA     = "DATA"
	EXTERN   = "EXTERN"
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TYPE     = "TYPE"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	BLANK    = "BLANK"
)

var keywords = map[string]TokenType{
	"module": MODULE,
	"import": IMPORT,
	"enum":   ENUM,
	"data":   DATA,
	"extern": EXTERN,
	"func":   FUNCTION,
	"let":    LET,
	"type":   TYPE,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"_":      BLANK,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
