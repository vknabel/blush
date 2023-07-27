package lexer

import (
	"strings"

	"github.com/vknabel/lithia/token"
)

func (l *Lexer) parseLeadingDecorations() []token.DecorativeToken {
	var decors []token.DecorativeToken
	for {
		tok := l.parseDecorativeToken()
		if tok == nil {
			return decors
		}
		decors = append(decors, *tok)
	}
}

func (l *Lexer) parseDecorativeToken() *token.DecorativeToken {
	var tok token.DecorativeToken
	switch {
	case l.ch == '#': // COMMENT
		tok.Type = token.COMMENT
		tok.Literal = l.parseInlineComment()
	case l.ch == '/': // eventually COMMENT
		if l.peekChar() == '/' {
			tok.Type = token.COMMENT
			l.advance()
			tok.Literal = l.parseInlineComment()
		} else {
			return nil
		}
	case isWhitespace(l.ch):
		tok.Type = token.WHITESPACE
		tok.Literal = l.skipWhitespace()
	default:
		return nil
	}
	return &tok
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

func (l *Lexer) skipWhitespace() string {
	ws := strings.Builder{}
	for isWhitespace(l.ch) {
		ws.WriteByte(l.ch)
		l.advance()
	}
	return ws.String()
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
