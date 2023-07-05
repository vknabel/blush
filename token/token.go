package token

type TokenType string

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
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
