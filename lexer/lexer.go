package lexer

import (
	"github.com/vknabel/lithia/token"
)

type Lexer struct {
	module   string
	fileName string
	input    string // TODO: change to io.Reader or *bufio.Reader
	startPos int    // start position of this token
	peekPos  int    // current reading position in input (after current char)
	currPos  int    // current position in input (points to current char)
	ch       byte   // current char under examination
}

func New(module, fileName, input string) Lexer {
	l := Lexer{
		module:   module,
		fileName: fileName,
		input:    input,
	}
	l.advance()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() // TODO: newlines should be respected

	l.startPos = l.currPos

	switch l.ch {
	case '!': // BANG, NEQ
		if l.peekChar() == '=' {
			tok = token.Token{Type: token.NEQ, Literal: "!="}
			l.advance()
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}
	case '+': // PLUS
		tok = l.newToken(token.PLUS, l.ch)
	case '-': // MINUS, new ARROW
		if l.peekChar() == '>' {
			tok = token.Token{Type: token.ARROW, Literal: "->"}
			l.advance()
		} else {
			tok = l.newToken(token.MINUS, l.ch)
		}
	case '*': // ASTERISK
		tok = l.newToken(token.ASTERISK, l.ch)
	case '/': // SLASH, COMMENT
		if l.peekChar() == '/' {
			tok.Type = token.COMMENT
			l.advance()
			tok.Literal = l.parseInlineComment()
		} else {
			tok = l.newToken(token.SLASH, l.ch)
		}

	case '<': // LT, LTE
		if l.peekChar() == '=' {
			tok = token.Token{Type: token.LTE, Literal: "<="}
			l.advance()
		} else {
			tok = l.newToken(token.LT, l.ch)
		}
	case '>': // GT, GTE
		if l.peekChar() == '=' {
			tok = token.Token{Type: token.GTE, Literal: ">="}
			l.advance()
		} else {
			tok = l.newToken(token.GT, l.ch)
		}
	case '=': // ASSIGN, EQ, ARROW
		switch l.peekChar() {
		case '=':
			tok = token.Token{Type: token.EQ, Literal: "=="}
			l.advance()
		case '>':
			tok = token.Token{Type: token.ARROW, Literal: "=>"}
			l.advance()
		default:
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case '&': // AND
		if l.peekChar() == '&' {
			tok = token.Token{Type: token.AND, Literal: "&&"}
			l.advance()
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	case '|': // OR
		if l.peekChar() == '|' {
			tok = token.Token{Type: token.OR, Literal: "||"}
			l.advance()
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}

	case ':': // COLON
		tok = l.newToken(token.COLON, l.ch)
	case '.': // DOT
		tok = l.newToken(token.DOT, l.ch)
	case ',': // COMMA
		tok = l.newToken(token.COMMA, l.ch)
	case '(': // LPAREN
		tok = l.newToken(token.LPAREN, l.ch)
	case ')': // RPAREN
		tok = l.newToken(token.RPAREN, l.ch)
	case '{': // LBRACE
		tok = l.newToken(token.LBRACE, l.ch)
	case '}': // RBRACE
		tok = l.newToken(token.RBRACE, l.ch)
	case '[': // LBRACKET
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']': // RBRACKET
		tok = l.newToken(token.RBRACKET, l.ch)
	case '@': // AT
		tok = l.newToken(token.AT, l.ch)

	case '#': // COMMENT
		tok.Type = token.COMMENT
		tok.Literal = l.parseInlineComment()
	case '"': // STRING
		tok.Type = token.STRING
		tok.Literal = l.parseString()
	case 0: // EOF
		tok.Type = token.EOF
	default: // IDENT, INT, FLOAT
		if isLetter(l.ch) {
			tok.Literal = l.parseIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.parseNumber()
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.advance()
	return tok
}

func (l *Lexer) parseString() string {
	position := l.currPos + 1
	for {
		l.advance()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.currPos]
}

func (l *Lexer) parseInlineComment() string {
	position := l.currPos + 1
	for {
		l.advance()
		if l.ch == '\n' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.currPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.advance()
	}
}

func (l *Lexer) advance() {
	if l.peekPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.peekPos]
	}
	l.currPos = l.peekPos
	l.peekPos += 1
}

func (l *Lexer) peekChar() byte {
	if l.peekPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPos]
	}
}

func (l *Lexer) parseIdentifier() string {
	position := l.currPos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.advance()
	}
	return l.input[position:l.currPos]
}

func (l *Lexer) parseNumber() (string, token.TokenType) {
	position := l.currPos
	for isDigit(l.ch) {
		l.advance()
	}
	if l.ch != '.' {
		return l.input[position:l.currPos], token.INT
	}
	if !isDigit(l.peekChar()) {
		return l.input[position:l.currPos], token.INT
	}
	l.advance()
	for isDigit(l.ch) {
		l.advance()
	}
	return l.input[position:l.currPos], token.FLOAT
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Source:  token.MakeSource(l.module, l.fileName, l.currPos),
	}
}
