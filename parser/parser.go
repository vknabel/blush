package parser

import (
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
		stmt := p.parseGlobalStatement()
		if stmt != nil {
			srcFile.Add(stmt)
		}
		p.nextToken()
	}
	return srcFile
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) peekIs(tokTypes ...token.TokenType) bool {
	for _, tok := range tokTypes {
		if p.peekToken.Type == tok {
			return true
		}
	}
	return false
}

func (p *Parser) expectPeek(tokTypes ...token.TokenType) bool {
	if !p.peekIs(tokTypes...) {
		p.detectError(UnexpectedGot(p.peekToken, tokTypes...))
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) detectError(err ParseError) {
	p.errors = append(p.errors, err)
}

func (p *Parser) parseGlobalStatement() ast.Statement {
	switch p.curToken.Type {
	case token.ENUM:
		return p.parseEnumDecl()
	case token.DATA:
		return p.parseDataDecl()
	case token.MODULE:
		return p.parseModuleDecl()
	case token.EXTERN:
		return p.parseExternDecl()
	case token.FUNCTION:
		return p.parseFunctionDecl()
	case token.IMPORT:
		return p.parseImportDecl()
	case token.AT:
		return p.parseAnnotationChainStatementDeclaration()
	case token.LET:
		return p.parseVariableDecl()
	default:
		return nil
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
func (p *Parser) parseEnumDecl() *ast.DeclEnum {
	return nil
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
func (p *Parser) parseDataDecl() *ast.DeclData {
	declToken := p.curToken
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	ident := ast.MakeIdentifier(p.curToken)
	data := ast.MakeDeclData(declToken, ident)

	if !p.peekIs(token.LBRACE) {
		return data
	}

	p.nextToken()

	// for {
	// TODO: allowed tokens:
	// - identifiers for members
	// - func and let declarations
	// - annotations preceding decls
	//}

	return data
}

func (p *Parser) parseModuleDecl() *ast.DeclModule {
	return nil
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
func (p *Parser) parseExternDecl() ast.StatementDeclaration {
	return nil
}

func (p *Parser) parseFunctionDecl() *ast.DeclFunc {
	return nil
}

func (p *Parser) parseImportDecl() *ast.DeclImport {
	return nil
}

// parseAnnotationChainStatementDeclaration parses annotations
//
// @ <fully-qualified-identifier>
// @ <fully-qualified-identifier>()
// @ <fully-qualified-identifier>(<param_list>)
func (p *Parser) parseAnnotationChainStatementDeclaration() ast.StatementDeclaration {
	return nil
}

func (p *Parser) parseVariableDecl() *ast.DeclVariable {
	return nil
}

func (p *Parser) parsePropertyDeclarationList() { // TODO: return?
	// TODO: allowed tokens:
	// - identifiers for members
	// - func and let declarations
	// - annotations preceding decls
	annotations := p.parseAnnotationChain()

	var field *ast.DeclField

	switch p.curToken.Type {
	case token.IDENT:
		var params []ast.DeclParameter
		if p.peekIs(token.LPAREN) {
			p.nextToken()
			params = p.parseParamList()
		}
		field = ast.MakeDeclField(ast.MakeIdentifier(p.curToken), params, annotations)
	case token.FUNCTION:
	case token.LET:
	}
	_ = field
}

func (p *Parser) parseAnnotationChain() *ast.AnnotationChain { // TODO: return?
	return nil
}

func (p *Parser) parseParamList() []ast.DeclParameter {
	params := make([]ast.DeclParameter, 0)

	for {
		annos := p.parseAnnotationChain()

		if !p.peekIs(token.IDENT) {
			// eventual errors will be triggered by parent
			return params
		}
		ident := ast.MakeIdentifier(p.curToken)
		params = append(params, *ast.MakeDeclParameter(ident, annos))
	}
}
