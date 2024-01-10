package parser

import (
	"fmt"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/lexer"
	"github.com/vknabel/lithia/token"
)

type Parser struct {
	lex    lexer.Lexer
	errors []ParseError

	curToken  token.Token
	peekToken token.Token
}

func New(lex lexer.Lexer) *Parser {
	p := &Parser{lex: lex}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) ParseSourceFile(filePath, moduleName string, input string) *ast.SourceFile {
	srcFile := ast.MakeSourceFile(filePath)

	for p.curToken.Type != token.EOF {
		stmt, decls := p.parseGlobalStatement()
		if stmt != nil {
			srcFile.Add(stmt)
			for _, d := range decls {
				srcFile.AddDecl(d)
			}
		} else {
			p.detectError(UnexpectedGot(p.curToken))
			p.nextToken()
		}
		// TODO: prevent infinite loop!
	}
	return srcFile
}

func (p *Parser) nextToken() {
	prev := fmt.Sprintf("%s(%s)", p.curToken.Type, p.peekToken.Type)
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
	fmt.Printf("%s -> %s(%s)\n", prev, p.curToken.Type, p.peekToken.Type)
}

func (p *Parser) peekIs(tokTypes ...token.TokenType) bool {
	for _, tok := range tokTypes {
		if p.peekToken.Type == tok {
			return true
		}
	}
	return false
}

func (p *Parser) curIs(tokTypes ...token.TokenType) bool {
	for _, tok := range tokTypes {
		if p.curToken.Type == tok {
			return true
		}
	}
	return false
}

func (p *Parser) inlinePeekIs(tokTypes ...token.TokenType) bool {
	for _, deco := range p.peekToken.Leading {
		if deco.Type != token.DECO_INLINE {
			return false
		}
	}
	return p.peekIs(tokTypes...)
}

func (p *Parser) expectCurToken(tokTypes ...token.TokenType) (token.Token, bool) {
	if !p.curIs(tokTypes...) {
		p.detectError(UnexpectedGot(p.peekToken, tokTypes...))
		return p.errorToken(), false
	}
	cur := p.curToken
	p.nextToken()
	return cur, true
}

func (p *Parser) expectPeekToken(tokTypes ...token.TokenType) (token.Token, bool) {
	if !p.peekIs(tokTypes...) {
		p.detectError(UnexpectedGot(p.peekToken, tokTypes...))
		return p.errorToken(), false
	}
	p.nextToken()
	return p.curToken, true
}

func (p *Parser) errorToken() token.Token {
	return token.Token{
		Type:    token.ILLEGAL,
		Literal: "ERROR",
		Source:  p.curToken.Source,
		Leading: p.curToken.Leading,
	}
}

func (p *Parser) detectError(err ParseError) {
	panic(err.Error())
	p.errors = append(p.errors, err)
}

func (p *Parser) parseGlobalStatement() (ast.Statement, []ast.Decl) {
	switch p.curToken.Type {
	case token.ENUM:
		return p.parseEnumDecl(nil)
	case token.DATA:
		return p.parseDataDecl(nil), nil
	case token.MODULE:
		return p.parseModuleDecl(nil), nil
	case token.EXTERN:
		return p.parseExternDecl(nil), nil
	case token.FUNCTION:
		return p.parseFunctionDecl(nil), nil
	case token.LET:
		return p.parseVariableDecl(nil), nil
	case token.IMPORT:
		return p.parseImportDecl(), nil
	case token.AT:
		return p.parseStatementDeclaration()
	default:
		err := UnexpectedGot(p.curToken, token.ENUM, token.DATA, token.MODULE, token.EXTERN, token.FUNCTION, token.IMPORT, token.AT, token.LET, token.IDENT, token.BLANK, token.IF, token.STRING, token.FLOAT, token.INT, token.LPAREN, token.LBRACKET, token.FOR, token.BANG)
		p.detectError(err)
		return nil, nil
	}
}

func (p *Parser) parseStatementDeclaration() (ast.StatementDeclaration, []ast.Decl) {
	annos := p.parseAnnotationChain()
	switch p.curToken.Type {
	case token.ENUM:
		return p.parseEnumDecl(annos)
	case token.DATA:
		return p.parseDataDecl(annos), nil
	case token.MODULE:
		return p.parseModuleDecl(annos), nil
	case token.EXTERN:
		return p.parseExternDecl(annos), nil
	case token.FUNCTION:
		return p.parseFunctionDecl(annos), nil
	case token.LET:
		return p.parseVariableDecl(annos), nil
	default:
		err := UnexpectedGot(p.curToken, token.ENUM, token.DATA, token.MODULE, token.EXTERN, token.FUNCTION, token.LET)
		p.detectError(err)
		return nil, nil
	}
}

// parseEnumDecl parsed enum declarations in various forms:
//
//		enum <identifier> // empty enum
//		enum <identifier> { } // empty enum
//		enum <identifier> {
//		  <identifier> // referencing: no annotations allowed!
//		  <fully-qualified-identifier> // global reference
//		  <optional:annotations> <data_decl>
//		  <optional:annotations> <enum_decl>
//	 	}
func (p *Parser) parseEnumDecl(annos *ast.AnnotationChain) (*ast.DeclEnum, []ast.Decl) {
	enumToken, _ := p.expectPeekToken(token.ENUM)
	if _, ok := p.expectPeekToken(token.IDENT); !ok {
		return nil, nil
	}
	ident := ast.MakeIdentifier(p.curToken)
	enum := ast.MakeDeclEnum(enumToken, ident)
	enum.Annotations = annos

	if !p.peekIs(token.LBRACE) {
		return enum, nil
	}

	p.nextToken()

	var childDecls []ast.Decl
	for !p.peekIs(token.RBRACE) {
		enumCase, children := p.parseEnumDeclCase()
		childDecls = append(childDecls, children...)
		enum.AddCase(enumCase)
	}
	p.nextToken()

	return enum, childDecls
}

// parseEnumDeclCase parses enum cases in these forms:
//
//	<identifier> // referencing: no annotations allowed!
//	<fully-qualified-identifier> // global reference
//	<optional:annotations> <data_decl>
//	<optional:annotations> <enum_decl>
func (p *Parser) parseEnumDeclCase() (*ast.DeclEnumCase, []ast.Decl) {
	if p.curToken.Type == token.IDENT {
		ref := p.parseStaticIdentifierReference()
		return ast.MakeDeclEnumCase(ref.TokenLiteral(), ref), nil
	}
	annotations := p.parseOptionalAnnotationChain()
	switch p.curToken.Type {
	case token.DATA:
		dataDecl := p.parseDataDecl(annotations)
		return ast.MakeDeclEnumCase(dataDecl.Token, ast.StaticReference{dataDecl.DeclName()}), []ast.Decl{dataDecl}
	case token.ENUM:
		enumDecl, childDecls := p.parseEnumDecl(annotations)
		return ast.MakeDeclEnumCase(enumDecl.Token, ast.StaticReference{enumDecl.DeclName()}), append(childDecls, enumDecl)
	default:
		p.detectError(UnexpectedGot(p.curToken, token.DATA, token.ENUM))
		return nil, nil
	}
}

func (p *Parser) parseStaticIdentifierReference() ast.StaticReference {
	var ref ast.StaticReference
	for p.peekIs(token.DOT) {
		id := ast.MakeIdentifier(p.curToken)
		ref = append(ref, id)
		p.nextToken()
	}
	if len(ref) == 0 {
		if _, ok := p.expectPeekToken(token.IDENT); ok {
			panic("broken invariant")
		}
		return ast.StaticReference{ast.MakeIdentifier(p.errorToken())}
	}
	return ref
}

// parseDataDecl parses data declarations in various forms:
//
//	data <identifier> // name
//	data <identifier> { }
//	data <identifier> {
//	  <optional:annotations> <identifer> // property name
//	  <optional:annotations> <identifier>(<param_list>) // function member
//	  <optional:annotations> <identifier> = <expr> // defaulted member
//
//	  // optional
//	  <optional:annotations> <func_decl>
//	  <optional:annotations> <var_decl>
//	}
func (p *Parser) parseDataDecl(annos *ast.AnnotationChain) *ast.DeclData {
	declToken, _ := p.expectCurToken(token.DATA)
	identToken, _ := p.expectCurToken(token.IDENT)
	ident := ast.MakeIdentifier(identToken)
	data := ast.MakeDeclData(declToken, ident)
	data.Annotations = annos

	if !p.curIs(token.LBRACE) {
		return data
	}
	p.expectCurToken(token.LBRACE)
	fields := p.parsePropertyDeclarationList()
	for _, f := range fields {
		data.AddField(f)
	}
	p.expectCurToken(token.RBRACE)
	return data
}

func (p *Parser) parseDataDeclField() *ast.DeclField {
	annotations := p.parseOptionalAnnotationChain()
	if p.curToken.Type != token.IDENT {
		p.detectError(UnexpectedGot(p.curToken, token.IDENT))
		return nil
	}
	name := ast.MakeIdentifier(p.curToken)
	p.nextToken()

	if !p.peekIs(token.LPAREN) {
		return ast.MakeDeclField(name, nil, annotations)
	}

	p.nextToken()

	params := p.parseParamList()

	p.expectPeekToken(token.RPAREN)
	return ast.MakeDeclField(name, params, annotations)
}

func (p *Parser) parseModuleDecl(annos *ast.AnnotationChain) *ast.DeclModule {
	modToken := p.curToken
	p.expectPeekToken(token.MODULE)

	name := ast.MakeIdentifier(p.curToken)
	p.nextToken()

	mod := ast.MakeDeclModule(modToken, name)
	mod.Annotations = annos
	return mod
}

// parseExternDecl parses two possible types:
// 1. an external function
// 2. an external type
//
//	extern <identifier> // type
//	extern <identifier>() // function
//	extern <identifier> { // type
//	  // see data declarations
//	}
func (p *Parser) parseExternDecl(annos *ast.AnnotationChain) ast.StatementDeclaration {
	externTok, _ := p.expectPeekToken(token.EXTERN)
	nameTok, _ := p.expectPeekToken(token.IDENT)
	nameIdent := ast.MakeIdentifier(nameTok)

	if p.peekIs(token.LBRACE) {
		// TODO: parse properties
		p.expectPeekToken(token.LBRACE)
		extern := ast.MakeDeclExternType(externTok, nameIdent)
		extern.Annotations = annos
		fields := p.parsePropertyDeclarationList()
		for _, f := range fields {
			extern.AddField(f)
		}
		p.expectPeekToken(token.RBRACE)
		return extern
	}
	var params []ast.DeclParameter
	if p.peekIs(token.LPAREN) {
		p.expectPeekToken(token.LPAREN)
		params = p.parseParamList()
		p.expectPeekToken(token.RPAREN)
	}
	extern := ast.MakeDeclExternFunc(externTok, nameIdent, params)
	extern.Annotations = annos
	return extern
}

func (p *Parser) parseFunctionDecl(annos *ast.AnnotationChain) *ast.DeclFunc {
	return nil
}

func (p *Parser) parseImportDecl() *ast.DeclImport {
	return nil
}

func (p *Parser) parseVariableDecl(annos *ast.AnnotationChain) *ast.DeclVariable {
	return nil
}

func (p *Parser) parsePropertyDeclarationList() []ast.DeclField {
	var fields []ast.DeclField
	for {
		// p.nextToken()
		if p.curToken.Type == token.RBRACE {
			return fields
		}
		field := p.parseDataDeclField()
		if field != nil {
			fields = append(fields, *field)
		}
	}
}

func (p *Parser) parseOptionalAnnotationChain() *ast.AnnotationChain { // TODO: return?
	if p.curToken.Type != token.AT {
		return nil
	}
	return p.parseAnnotationChain()
}

func (p *Parser) parseAnnotationChain() *ast.AnnotationChain {
	atTok := p.curToken
	p.expectPeekToken(token.AT)

	ref := p.parseStaticIdentifierReference()

	chain := ast.MakeAnnotationChain(atTok, ref)

	if p.peekToken.Type != token.LPAREN {
		return chain
	}
	p.expectPeekToken(token.LPAREN)
	// TODO: parameters
	p.expectPeekToken(token.RPAREN)
	return chain
}

func (p *Parser) parseParamList() []ast.DeclParameter {
	params := make([]ast.DeclParameter, 0)

	for {
		annos := p.parseOptionalAnnotationChain()

		if !p.peekIs(token.IDENT) {
			// eventual errors will be triggered by parent
			return params
		}
		p.nextToken()
		ident := ast.MakeIdentifier(p.curToken)
		params = append(params, *ast.MakeDeclParameter(ident, annos))
	}
}
