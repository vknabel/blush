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

	prefixParsers map[token.TokenType]prefixParser
	infixParsers  map[token.TokenType]infixParser
}

func New(lex lexer.Lexer) *Parser {
	p := &Parser{lex: lex}
	p.nextToken()
	p.nextToken()

	p.prefixParsers = make(map[token.TokenType]prefixParser)
	p.registerPrefix(token.IDENT, p.parsePrattExprIdentifier)
	p.registerPrefix(token.MINUS, p.parsePrattExprPrefix)
	p.registerPrefix(token.INT, p.parsePrattExprInt)
	p.registerPrefix(token.FLOAT, p.parsePrattExprFloat)
	p.registerPrefix(token.BANG, p.parsePrattExprPrefix)
	p.registerPrefix(token.MINUS, p.parsePrattExprPrefix)
	p.registerPrefix(token.LPAREN, p.parsePrattExprGroup)
	p.registerPrefix(token.IF, p.parsePrattExprIfElse) // only exactly one expr per if / else if / else, else mandatory, later we eventually want to allow assignments and local vars
	p.registerPrefix(token.LBRACE, p.parsePrattExprFunc)
	// p.registerPrefix(token.TYPE, p.parseExprType) // only exactly one expr per case
	// p.registerPrefix(token.SWITCH / MATCH, p.parseExprSwitch) // only exactly one expr per case
	p.registerPrefix(token.LBRACKET, p.parseExprListOrDict)
	p.registerPrefix(token.STRING, p.parsePrattExprString)

	p.infixParsers = make(map[token.TokenType]infixParser)
	p.registerInfix(token.EQ, p.parsePrattExprInfix)
	p.registerInfix(token.NEQ, p.parsePrattExprInfix)
	p.registerInfix(token.LTE, p.parsePrattExprInfix)
	p.registerInfix(token.GTE, p.parsePrattExprInfix)
	p.registerInfix(token.LT, p.parsePrattExprInfix)
	p.registerInfix(token.GT, p.parsePrattExprInfix)
	p.registerInfix(token.PLUS, p.parsePrattExprInfix)
	p.registerInfix(token.MINUS, p.parsePrattExprInfix)
	p.registerInfix(token.SLASH, p.parsePrattExprInfix)
	p.registerInfix(token.ASTERISK, p.parsePrattExprInfix)
	p.registerInfix(token.LPAREN, p.parsePrattExprCall)
	p.registerInfix(token.DOT, p.parsePrattExprMember)
	p.registerInfix(token.LBRACKET, p.parsePrattExprIndex)

	return p
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) ParseSourceFile(filePath, moduleName string, input string) *ast.SourceFile {
	srcFile := ast.MakeSourceFile(filePath)

	for p.curToken.Type != token.EOF {
		stmt, childDecls := p.parseGlobalStatement()
		if stmt != nil {
			srcFile.Add(stmt)
			for _, d := range childDecls {
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

func (p *Parser) nextToken() token.Token {
	cur := p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
	return cur
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

func (p *Parser) expect(tokTypes ...token.TokenType) (token.Token, bool) {
	if !p.curIs(tokTypes...) {
		p.detectError(UnexpectedGot(p.curToken, tokTypes...))
		return p.errorToken(), false
	}
	cur := p.curToken
	p.nextToken()
	return cur, true
}

func (p *Parser) skip(tokTypes ...token.TokenType) {
	if !p.curIs(tokTypes...) {
		return
	}
	p.nextToken()
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
	case token.ANNOTATION:
		return p.parseAnnotationDecl(nil), nil
	case token.EXTERN:
		return p.parseExternDecl(nil), nil
	case token.FUNCTION:
		return p.parseFunctionDecl(nil), nil
	case token.LET:
		return p.parseVariableDecl(nil), nil
	case token.IMPORT:
		return p.parseImportDecl(), nil
	case token.AT:
		return p.parseAnnotatedStatementDeclaration()
	default:
		if _, ok := p.prefixParsers[p.curToken.Type]; ok {
			return p.parseExprStmt(), nil
		}

		prefixes := []token.TokenType{
			token.ENUM, token.DATA, token.MODULE, token.EXTERN, token.FUNCTION, token.IMPORT, token.AT, token.LET, token.IF, token.FOR,
		}
		for t := range p.prefixParsers {
			prefixes = append(prefixes, t)
		}
		err := UnexpectedGot(p.curToken, prefixes...)
		p.detectError(err)
		return nil, nil
	}
}

func (p *Parser) parseAnnotatedStatementDeclaration() (ast.StatementDeclaration, []ast.Decl) {
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
func (p *Parser) parseEnumDecl(annos ast.AnnotationChain) (*ast.DeclEnum, []ast.Decl) {
	enumToken, _ := p.expect(token.ENUM)
	identToken, _ := p.expect(token.IDENT)
	ident := ast.MakeIdentifier(identToken)
	enum := ast.MakeDeclEnum(enumToken, ident)
	enum.Annotations = annos

	if !p.curIs(token.LBRACE) {
		return enum, nil
	}

	p.expect(token.LBRACE)

	var childDecls []ast.Decl
	for !p.curIs(token.RBRACE) {
		enumCase, children := p.parseEnumDeclCase()
		childDecls = append(childDecls, children...)
		enum.AddCase(enumCase)
	}
	p.expect(token.RBRACE)

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
	annotations := p.parseAnnotationChain()
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
	for {
		identTok, ok := p.expect(token.IDENT)
		if !ok {
			break
		}
		id := ast.MakeIdentifier(identTok)
		ref = append(ref, id)

		if !p.curIs(token.DOT) {
			break
		}
		p.expect(token.DOT)
	}
	if len(ref) == 0 {
		if _, ok := p.expectPeekToken(token.IDENT); ok {
			panic(fmt.Sprintf("invariant error: empty static reference, token:%+v", p.curToken))
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
func (p *Parser) parseDataDecl(annos ast.AnnotationChain) *ast.DeclData {
	declToken, _ := p.expect(token.DATA)
	identToken, _ := p.expect(token.IDENT)
	ident := ast.MakeIdentifier(identToken)
	data := ast.MakeDeclData(declToken, ident)
	data.Annotations = annos

	if !p.curIs(token.LBRACE) {
		return data
	}
	p.expect(token.LBRACE)
	fields := p.parsePropertyDeclarationList()
	for _, f := range fields {
		data.AddField(f)
	}
	p.expect(token.RBRACE)
	return data
}

func (p *Parser) parseDataDeclField() *ast.DeclField {
	annotations := p.parseAnnotationChain()
	identTok, _ := p.expect(token.IDENT)
	name := ast.MakeIdentifier(identTok)

	if !p.curIs(token.LPAREN) {
		return ast.MakeDeclField(name, nil, annotations)
	}

	p.expect(token.LPAREN)
	params := p.parseDeclParameterList()
	p.expect(token.RPAREN)
	return ast.MakeDeclField(name, params, annotations)
}

func (p *Parser) parseAnnotationDecl(annos *ast.AnnotationChain) *ast.DeclAnnotation {
	declToken, _ := p.expect(token.ANNOTATION)
	identToken, _ := p.expect(token.IDENT)
	ident := ast.MakeIdentifier(identToken)
	declAnno := ast.MakeDeclAnnotation(declToken, ident)
	declAnno.Annotations = annos

	if !p.curIs(token.LBRACE) {
		return declAnno
	}
	p.expect(token.LBRACE)
	fields := p.parsePropertyDeclarationList()
	for _, f := range fields {
		declAnno.AddField(f)
	}
	p.expect(token.RBRACE)
	return declAnno
}

func (p *Parser) parseModuleDecl(annos ast.AnnotationChain) *ast.DeclModule {
	modToken, _ := p.expect(token.MODULE)
	nameTok, _ := p.expect(token.IDENT)
	name := ast.MakeIdentifier(nameTok)

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
func (p *Parser) parseExternDecl(annos ast.AnnotationChain) ast.StatementDeclaration {
	externTok, _ := p.expect(token.EXTERN)
	nameTok, _ := p.expect(token.IDENT)
	nameIdent := ast.MakeIdentifier(nameTok)

	if p.curIs(token.LBRACE) {
		// TODO: parse properties
		p.expect(token.LBRACE)
		extern := ast.MakeDeclExternType(externTok, nameIdent)
		extern.Annotations = annos
		fields := p.parsePropertyDeclarationList()
		for _, f := range fields {
			extern.AddField(f)
		}
		p.expect(token.RBRACE)
		return extern
	}
	var params []ast.DeclParameter
	if p.curIs(token.LPAREN) {
		p.expect(token.LPAREN)
		params = p.parseDeclParameterList()
		p.expect(token.RPAREN)
	}
	extern := ast.MakeDeclExternFunc(externTok, nameIdent, params)
	extern.Annotations = annos
	return extern
}

func (p *Parser) parseFunctionDecl(annos ast.AnnotationChain) *ast.DeclFunc {
	funcTok, _ := p.expect(token.FUNCTION)
	nameTok, _ := p.expect(token.IDENT)

	var impl *ast.ExprFunc
	if p.curIs(token.LPAREN) {
		p.expect(token.LPAREN)
		params := p.parseDeclParameterList()
		p.expect(token.RPAREN)
		implTok, _ := p.expect(token.LBRACE)
		block := p.parseStmtBlock()
		p.expect(token.RBRACE)
		impl = ast.MakeExprFunc(implTok, nameTok.Literal, params, block)
	} else {
		impl = p.parseExprFunction()
	}

	decl := ast.MakeDeclFunc(funcTok, ast.MakeIdentifier(nameTok), impl)
	decl.Annotations = annos
	return decl
}

func (p *Parser) parseImportDecl() *ast.DeclImport {
	importTok, _ := p.expect(token.IMPORT)

	var importDecl *ast.DeclImport
	if p.peekIs(token.ASSIGN) {
		aliasTok, _ := p.expect(token.IDENT)
		p.expect(token.ASSIGN)
		moduleName := p.parseStaticIdentifierReference()
		importDecl = ast.MakeDeclAliasImport(importTok, ast.MakeIdentifier(aliasTok), moduleName)
	} else {
		moduleName := p.parseStaticIdentifierReference()
		importDecl = ast.MakeDeclImport(importTok, moduleName)
	}

	if !p.curIs(token.LBRACE) {
		return importDecl
	}
	p.expect(token.LBRACE)
	for !p.curIs(token.RBRACE) {
		memberTok, _ := p.expect(token.IDENT)
		member := ast.MakeDeclImportMember(memberTok, importDecl.ModuleName, ast.MakeIdentifier(memberTok))
		importDecl.AddMember(member)

		if p.curIs(token.COMMA) {
			p.expect(token.COMMA)
		}
	}
	p.expect(token.RBRACE)
	return importDecl
}

func (p *Parser) parseVariableDecl(annos ast.AnnotationChain) *ast.DeclVariable {
	letTok, _ := p.expect(token.LET)
	nameTok, _ := p.expect(token.IDENT)
	name := ast.MakeIdentifier(nameTok)
	p.expect(token.ASSIGN)
	expr := p.parseExpr()
	let := ast.MakeDeclVariable(letTok, name, expr)
	let.Annotations = annos
	return let
}

func (p *Parser) parsePropertyDeclarationList() []ast.DeclField {
	var fields []ast.DeclField
	for {
		if p.curToken.Type == token.RBRACE {
			return fields
		}
		field := p.parseDataDeclField()
		if field != nil {
			fields = append(fields, *field)
		}
	}
}

func (p *Parser) parseAnnotationChain() ast.AnnotationChain {
	var annotationChain ast.AnnotationChain
	for p.curIs(token.AT) {
		anno := p.parseAnnotationInstance()
		annotationChain = append(annotationChain, anno)
	}
	return annotationChain
}

func (p *Parser) parseAnnotationInstance() *ast.DeclAnnotationInstance {
	atTok, _ := p.expect(token.AT)
	ref := p.parseStaticIdentifierReference()

	anno := ast.MakeAnnotationInstance(atTok, ref)

	if !p.curIs(token.LPAREN) {
		return anno
	}
	p.expect(token.LPAREN)
	args := p.parseExprArgumentList()
	for _, arg := range args {
		anno.AddArgument(arg)
	}
	p.expect(token.RPAREN)
	return anno
}

func (p *Parser) parseDeclParameterList() []ast.DeclParameter {
	params := make([]ast.DeclParameter, 0)

	for {
		annos := p.parseAnnotationChain()
		if !p.curIs(token.IDENT) {
			// eventual errors will be triggered by parent
			return params
		}
		identTok, _ := p.expect(token.IDENT)
		ident := ast.MakeIdentifier(identTok)
		params = append(params, *ast.MakeDeclParameter(ident, annos))

		if !p.curIs(token.COMMA) {
			return params
		}
		p.expect(token.COMMA)
	}
}

func (p *Parser) parseStatementReturn() ast.StmtReturn {
	retTok, _ := p.expect(token.RETURN)
	for _, dec := range p.curToken.Leading {
		if dec.Type != token.DECO_INLINE {
			return ast.MakeStmtReturn(retTok, nil)
		}
	}
	expr := p.parseExpr()
	return ast.MakeStmtReturn(retTok, expr)
}

func (p *Parser) parseStatementIf() ast.StmtIf {
	ifTok, _ := p.expect(token.IF)
	cond := p.parseExpr()
	p.expect(token.LBRACE)
	ifBlock := p.parseStmtBlock()
	p.expect(token.RBRACE)

	ifStmt := ast.MakeStmtIf(ifTok, cond, ifBlock)

	for p.curIs(token.ELSE) {
		if p.peekIs(token.IF) {
			elseIf := p.parseStatementElseIf()
			ifStmt.AddElseIf(elseIf)
			continue
		}
		p.expect(token.ELSE)
		p.expect(token.LBRACE)
		elseBlock := p.parseStmtBlock()
		p.expect(token.RBRACE)
		ifStmt.SetElse(elseBlock)
		break
	}
	return ifStmt
}

func (p *Parser) parseStatementElseIf() ast.StmtElseIf {
	elseTok, _ := p.expect(token.ELSE)
	p.expect(token.IF)
	cond := p.parseExpr()
	p.expect(token.LBRACE)
	block := p.parseStmtBlock()
	p.expect(token.RBRACE)

	return ast.MakeStmtIfElse(elseTok, cond, block)
}

func (p *Parser) parseExprArgumentList() []ast.Expr {
	var args []ast.Expr
	for !p.curIs(token.RPAREN) {
		args = append(args, p.parseExpr())
		if !p.curIs(token.COMMA) {
			return args
		}
		p.expect(token.COMMA)
	}
	return args
}

func (p *Parser) parseExpr() ast.Expr {
	expr := p.parsePrattExpr(LOWEST)
	return expr
}

func (p *Parser) parseStmtBlock() ast.Block {
	block := make([]ast.Statement, 0)
	// TODO: parse statements and local decls
	return block
}

func (p *Parser) parseExprFunction() *ast.ExprFunc {
	tok, _ := p.expect(token.LBRACE)
	params := p.parseDeclParameterList()
	if len(params) == 0 {
		p.skip(token.ARROW)
	} else {
		p.expect(token.ARROW)
	}
	impl := p.parseStmtBlock()
	p.expect(token.RBRACE)
	return ast.MakeExprFunc(tok, "TODO", params, impl)
}
