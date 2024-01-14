package parser

import (
	"fmt"
	"strconv"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/token"
)

type (
	prefixParser func() ast.Expr
	infixParser  func(ast.Expr) ast.Expr

	Precedence    int
	Associativity int
)

const (
	_ Precedence = iota
	LOWEST
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	COMPARISON  // == or != or <= or >= or < or >
	COALESCING  // placeholder for ??
	RANGE       // placeholder for ..<
	SUM         // + or -
	PRODUCT     // * or / or %
	BITWISE     // placeholder for << and >>
	PREFIX      // -x or !x
	CALL        // fun(x)
	MEMBER      // . or ?.
)

var precedences = map[token.TokenType]Precedence{
	token.EQ:       COMPARISON,
	token.NEQ:      COMPARISON,
	token.LTE:      COMPARISON,
	token.GTE:      COMPARISON,
	token.LT:       COMPARISON,
	token.GT:       COMPARISON,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.DOT:      MEMBER,
}

const (
	A_NONE Associativity = iota
	A_LEFT
	A_RIGHT
)

var associativities = map[Precedence]Associativity{
	LOWEST:      A_NONE,
	LOGICAL_OR:  A_LEFT,
	LOGICAL_AND: A_LEFT,
	COMPARISON:  A_NONE,
	COALESCING:  A_RIGHT,
	RANGE:       A_NONE,
	SUM:         A_LEFT,
	PRODUCT:     A_LEFT,
	BITWISE:     A_NONE,
}

func (p *Parser) peekPrecedence() Precedence {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecendence() Precedence {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParser) {
	p.prefixParsers[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParser) {
	p.infixParsers[tokenType] = fn
}

func (p *Parser) parseExprStmt() *ast.StmtExpr {
	stmtTok := p.curToken
	expr := p.parsePrattExpr(LOWEST)
	p.nextToken()
	return ast.MakeStmtExpr(stmtTok, expr)
}

func (p *Parser) parsePrattExpr(precedence Precedence) ast.Expr {
	prefix := p.prefixParsers[p.curToken.Type]
	if prefix == nil {
		expectedTypes := make([]token.TokenType, 0, len(p.prefixParsers))
		for t := range p.prefixParsers {
			expectedTypes = append(expectedTypes, t)
		}
		p.expect(expectedTypes...)
		return nil
	}
	lhs := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return lhs
		}
		p.nextToken()
		lhs = infix(lhs)
	}
	return lhs
}

func (p *Parser) parsePrattExprIdentifier() ast.Expr {
	return ast.MakeExprIdentifier(ast.MakeIdentifier(p.curToken))
}

func (p *Parser) parsePrattExprInt() ast.Expr {
	int, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.detectError(UnderlyingErr{p.curToken, err})
	}
	return ast.MakeExprInt(int, p.curToken)
}

func (p *Parser) parsePrattExprFloat() ast.Expr {
	float, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.detectError(UnderlyingErr{p.curToken, err})
	}
	return ast.MakeExprFloat(float, p.curToken)
}

func (p *Parser) parsePrattExprPrefix() ast.Expr {
	op := ast.OperatorUnary(p.curToken)
	p.nextToken()
	expr := p.parsePrattExpr(PREFIX)
	return ast.MakeExprOperatorUnary(op, expr)
}

func (p *Parser) parsePrattExprInfix(lhs ast.Expr) ast.Expr {
	op := ast.OperatorBinary(p.curToken)

	prec := p.curPrecendence()
	asso := associativities[p.curPrecendence()]
	if asso == A_RIGHT {
		prec -= 1
	}
	p.nextToken()
	rhs := p.parsePrattExpr(prec)
	return ast.MakeExprOperatorBinary(op, lhs, rhs)
}

func (p *Parser) parsePrattExprGroup() ast.Expr {
	p.expect(token.LPAREN)
	expr := p.parsePrattExpr(LOWEST)

	_, ok := p.expectPeekToken(token.RPAREN)
	if !ok {
		return nil
	}
	return expr
}

func (p *Parser) parsePrattExprIfElse() ast.Expr {
	ifTok := p.curToken
	p.nextToken()

	condition := p.parsePrattExpr(LOWEST)
	if condition == nil {
		return nil
	}

	_, ok := p.expectPeekToken(token.LBRACE)
	if !ok {
		return nil
	}
	p.nextToken()
	then := p.parsePrattExpr(LOWEST)
	if then == nil {
		return nil
	}

	_, ok = p.expectPeekToken(token.RBRACE)
	if !ok {
		return nil
	}

	elseTok, ok := p.expectPeekToken(token.ELSE)
	if !ok {
		return nil
	}

	ifExpr := ast.MakeExprIf(ifTok, condition, then)

	var i int
	for p.curIs(token.ELSE) && p.peekIs(token.IF) {
		i++
		fmt.Printf("%d:%d %q -> %q\n", i, 1, p.curToken.Type, p.peekToken.Type)
		p.nextToken()
		p.nextToken()
		fmt.Printf("%d:%d %q -> %q\n", i, 2, p.curToken.Type, p.peekToken.Type)
		elseCond := p.parsePrattExpr(LOWEST)
		fmt.Printf("%T\n", elseCond)
		fmt.Printf("%d:%d %q -> %q\n", i, 3, p.curToken.Type, p.peekToken.Type)
		p.expectPeekToken(token.LBRACE)
		p.nextToken()
		fmt.Printf("%d:%d %q -> %q\n", i, 4, p.curToken.Type, p.peekToken.Type)
		elseExpr := p.parsePrattExpr(LOWEST)
		fmt.Printf("%d:%d %q -> %q\n", i, 5, p.curToken.Type, p.peekToken.Type)
		p.expectPeekToken(token.RBRACE)
		fmt.Printf("%d:%d %q -> %q\n", i, 6, p.curToken.Type, p.peekToken.Type)
		elif := ast.MakeExprElseIf(elseTok, elseCond, elseExpr)
		ifExpr.AddElseIf(elif)
		p.expectPeekToken(token.ELSE)
		fmt.Printf("%d:%d %q -> %q\n", i, 7, p.curToken.Type, p.peekToken.Type)
	}

	_, ok = p.expectPeekToken(token.LBRACE)
	if !ok {
		return nil
	}
	p.nextToken()
	els := p.parsePrattExpr(LOWEST)
	if els == nil {
		return nil
	}

	_, ok = p.expectPeekToken(token.RBRACE)
	if !ok {
		return nil
	}

	ifExpr.SetElse(els)
	return ifExpr
}

func (p *Parser) parsePrattExprFunc() ast.Expr {
	tok := p.curToken
	params := p.parseDeclParameterList()
	p.expect(token.ARROW)
	impl := p.parseStmtBlock()
	if !p.curIs(token.RBRACE) {
		p.expect(token.RBRACE)
	}
	return ast.MakeExprFunc(tok, "TODO", params, impl)
}

func (p *Parser) parsePrattExprCall(fn ast.Expr) ast.Expr {
	fnExpr := ast.MakeExprInvocation(fn)

	if p.peekIs(token.RPAREN) {
		p.nextToken()
		return fnExpr
	}
	p.nextToken()
	fnExpr.AddArgument(p.parsePrattExpr(LOWEST))

	for p.peekIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		fnExpr.AddArgument(p.parsePrattExpr(LOWEST))
	}

	_, ok := p.expectPeekToken(token.RPAREN)
	if !ok {
		return nil
	}

	return fnExpr
}

func (p *Parser) parsePrattExprMember(owner ast.Expr) ast.Expr {
	dotTok := p.curToken
	identTok, ok := p.expectPeekToken(token.IDENT)
	if !ok {
		return nil
	}
	return ast.MakeExprMemberAccess(dotTok, owner, ast.MakeIdentifier(identTok))
}

func (p *Parser) parsePrattExprIndex(owner ast.Expr) ast.Expr {
	indexTok := p.curToken
	indexExpr := p.parsePrattExpr(LOWEST)
	_, ok := p.expectPeekToken(token.RBRACKET)
	if !ok {
		return nil
	}
	return ast.MakeExprIndexAccess(indexTok, owner, indexExpr)
}

func (p *Parser) parsePrattExprString() ast.Expr {
	tok := p.curToken
	return ast.MakeExprString(tok, tok.Literal)
}

func (p *Parser) parseExprListOrDict() ast.Expr {
	tok := p.curToken

	if p.peekIs(token.RBRACKET) {
		p.nextToken()
		return ast.MakeExprArray(nil, tok)
	}
	if p.peekIs(token.COLON) {
		p.nextToken()
		return ast.MakeExprDict(nil, tok)
	}

	initialExpr := p.parsePrattExpr(LOWEST)

	if p.peekIs(token.COMMA) {
		rest := p.parsePrattExprArrayElements()
		if rest == nil {
			return nil
		}
		elements := append([]ast.Expr{initialExpr}, rest...)
		_, ok := p.expectPeekToken(token.RBRACKET)
		if !ok {
			return nil
		}
		return ast.MakeExprArray(elements, tok)
	}

	_, ok := p.expectPeekToken(token.COLON)
	if !ok {
		return nil
	}
	valueExpr := p.parsePrattExpr(LOWEST)
	initialEntry := ast.MakeExprDictEntry(initialExpr, valueExpr)
	entries := []ast.ExprDictEntry{initialEntry}

	if p.peekIs(token.RBRACKET) {
		p.nextToken()
		return ast.MakeExprDict(entries, tok)
	}
	rest := p.parsePrattExprDictEntries()
	if rest == nil {
		return nil
	}
	entries = append(entries, rest...)
	return ast.MakeExprDict(entries, tok)
}

func (p *Parser) parsePrattExprArrayElements() []ast.Expr {
	var elements []ast.Expr
	for p.peekIs(token.COMMA) {
		p.nextToken()

		expr := p.parsePrattExpr(LOWEST)
		if expr == nil {
			return nil
		}
		elements = append(elements, expr)
	}
	return elements
}

func (p *Parser) parsePrattExprDictEntries() []ast.ExprDictEntry {
	var elements []ast.ExprDictEntry
	for p.peekIs(token.COMMA) {
		p.nextToken()

		key := p.parsePrattExpr(LOWEST)
		if key == nil {
			return nil
		}
		_, ok := p.expectPeekToken(token.COLON)
		if !ok {
			return nil
		}
		value := p.parsePrattExpr(LOWEST)
		if value == nil {
			return nil
		}
		elements = append(elements, ast.MakeExprDictEntry(key, value))
	}
	return elements
}
