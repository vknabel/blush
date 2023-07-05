package parser

import (
	"github.com/vknabel/lithia/lexer"
	"github.com/vknabel/lithia/token"
)

type Parser struct {
	lex    *lexer.Lexer
	errors []ParseError

	curToken  token.Token
	peekToken token.Token
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

// func (p *Parser) ParseProgram() *ast.Program {
// 	panic("TODO")
// }

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}
