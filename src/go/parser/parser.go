// parser/parser.go

package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"z/ast"
	"z/lexer"
	"z/token"
)

var runSourceDir = ""

type (
	prefixParseFn func() ast.Expression
	infixPasrseFn func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	QUESTION
	LOGIC
	EQUALS
	LESSGRATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:             EQUALS,
	token.NOT_EQ:         EQUALS,
	token.AND:            LOGIC,
	token.OR:             LOGIC,
	token.LT:             LESSGRATER,
	token.GT:             LESSGRATER,
	token.GE:             LESSGRATER,
	token.LE:             LESSGRATER,
	token.PLUSASSIGN:     LESSGRATER,
	token.MINUSASSIGN:    LESSGRATER,
	token.ASTERISKASSIGN: LESSGRATER,
	token.SLASHASSIGN:    LESSGRATER,
	token.PLUSPLUS:       LESSGRATER,
	token.MINUSMINUS:     LESSGRATER,
	token.ASSIGN:         ASSIGN,
	token.CLASS:          LESSGRATER,
	token.PLUS:           SUM,
	token.MINUS:          SUM,
	token.SLASH:          PRODUCT,
	token.ASTERISK:       PRODUCT,
	token.LPAREN:         CALL,
	token.LBRACKET:       INDEX,
	token.OBJET_GET:      INDEX,
	token.CLASS_GET:      INDEX,
	token.QUESTION:       QUESTION,
}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixPasrseFns map[token.TokenType]infixPasrseFn

	tokenCount int
}

var initReadCount int = 2

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:          l,
		errors:     []string{},
		tokenCount: 0,
	}

	// 读取两个词法单元，以设置curToken和peekToken
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixPasrseFns = make(map[token.TokenType]infixPasrseFn)
	p.registerPrefix(token.IDENT, p.parseInditifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerInfix(token.PLUSASSIGN, p.parseInfixExpression)
	p.registerInfix(token.MINUSASSIGN, p.parseInfixExpression)
	p.registerInfix(token.ASTERISKASSIGN, p.parseInfixExpression)
	p.registerInfix(token.SLASHASSIGN, p.parseInfixExpression)
	p.registerInfix(token.PLUSPLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUSMINUS, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.CLASS, p.parseClassExpression)
	p.registerPrefix(token.INTERFACE, p.parseInterfaceExpress)
	p.registerPrefix(token.NEW, p.parseNewExpression)
	p.registerInfix(token.OBJET_GET, p.parseInfixExpression)
	p.registerInfix(token.CLASS_GET, p.parseInfixExpression)
	p.registerPrefix(token.DEFER, p.parseDeferExpression)
	p.registerInfix(token.QUESTION, p.parseQuestionExpression)
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	p.tokenCount = p.tokenCount + 1
}

func (p *Parser) nextNextToken() {
	p.nextToken()
	p.nextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		if p.curToken.Type == token.IMPORT {
			p.parseImportFile(program, p.peekToken.Literal)
		} else {
			stmt := p.parseStatement()
			if stmt != nil {
				program.Statements = append(program.Statements, stmt)
			}
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseImportFile(program *ast.Program, fileName string) {
	if runSourceDir != "" {
		if !strings.HasSuffix(fileName, ".z") {
			fileName = fileName + ".z"
		}
		if !strings.Contains(fileName, "builtin.z") {
			fileDir := filepath.Dir(p.l.FileName)
			fileDirArr := strings.Split(fileDir, runSourceDir)
			if len(fileDirArr) > 1 {
				fileName = runSourceDir + fileDirArr[1] + "/" + fileName
			} else {
				fileName = runSourceDir + "/" + fileName
			}
		}
	}
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("import file not exists: " + fileName)
			os.Exit(1)
		}
	}
	importCode, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	importLexer := lexer.New(string(importCode))
	importParser := New(importLexer)
	importLexer.SetFileName(fileName)
	importProgram := importParser.ParseProgram()
	program.Statements = append(program.Statements, importProgram.Statements...)
	p.nextToken() // remove file path string
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.PACKAGE:
		return p.parsePackageStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parsePackageStatement() ast.Statement {
	if p.tokenCount != initReadCount {
		p.errors = append(p.errors, "package need be the first token")
		return nil
	}
	p.nextToken()
	p.l.SetPackageName(p.curToken.Literal)
	p.nextToken()
	return nil
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken, FileName: p.l.FileName, PackageName: p.l.PackageName}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt.Name.FileName = p.l.FileName
	stmt.Name.PackageName = p.l.PackageName

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken, FileName: p.l.FileName, PackageName: p.l.PackageName}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stms := &ast.ExpressionStatement{Token: p.curToken, FileName: p.l.FileName, PackageName: p.l.PackageName}
	stms.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stms
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() && !p.peekTokenIs(token.COLON) {
		infix := p.infixPasrseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseInditifier() ast.Expression {
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifier.FileName = p.l.FileName
	identifier.PackageName = p.l.PackageName
	if p.peekTokenIs(token.LBRACKET) {
		return p.parseAssignHashExpress(identifier)
	}
	return identifier
}

func (p *Parser) parseAssignHashExpress(identifier *ast.Identifier) ast.Expression {
	p.nextToken()
	p.nextToken()
	index := p.parseExpression(LOWEST)
	p.nextToken()
	if !p.peekTokenIs(token.ASSIGN) {
		exp := &ast.IndexExpression{Token: identifier.Token, Left: identifier}
		exp.Index = index
		return exp
	}
	p.nextToken()
	p.nextToken()
	value := p.parseExpression(LOWEST)
	stmt := &ast.HashAssignExpress{Token: p.curToken, Hash: *identifier, Index: index, Value: value}
	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.RPAREN) {
		oneToken := token.Token{Literal: "1", Type: token.INT}
		expression.Right = ast.Expression(&ast.IntegerLiteral{Value: 1, Token: oneToken})
	} else {
		p.nextToken()
		expression.Right = p.parseExpression(precedence)
	}
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Body = p.parseBlockStatement()
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		expression, ok := stmt.(*ast.ExpressionStatement)
		isDeferBlock := false
		if ok {
			deferBlock, ok := expression.Expression.(*ast.BlockStatement)
			if ok {
				if deferBlock.IsDeferBlock {
					isDeferBlock = true
					block.DeferStatements = deferBlock.Statements
				}
			}
		}
		if stmt != nil && !isDeferBlock {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken, PackageName: p.l.PackageName, FileName: p.l.FileName}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextNextToken()
		ident.Default = p.parseExpression(LOWEST)
	}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextNextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if p.peekTokenIs(token.ASSIGN) {
			p.nextNextToken()
			ident.Default = p.parseExpression(LOWEST)
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken, Keys: make([]ast.Expression, 0)}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		hash.Keys = append(hash.Keys, key)
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) && !p.expectPeek(token.SEMICOLON) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixPasrseFn) {
	p.infixPasrseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) SetRunSourceDir(sourceDir string) {
	runSourceDir = sourceDir
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseBreakExpression() ast.Expression {
	return &ast.BreakExpression{Token: p.curToken}
}

func (p *Parser) parseForExpression() ast.Expression {
	forExpression := &ast.ForExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	forExpression.Initor = p.parseLetStatement()
	p.nextToken()
	forExpression.Condition = p.parseExpression(LOWEST)
	p.nextToken()
	p.nextToken()
	forExpression.After = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	forExpression.Body = p.parseBlockStatement()
	return forExpression
}

func (p *Parser) parseClassExpression() ast.Expression {
	classExpression := &ast.ClassExpress{
		Token: p.curToken,
	}
	p.expectPeek(token.IDENT)
	classExpression.Name = &ast.Identifier{
		Token:       p.curToken,
		Value:       p.curToken.Literal,
		PackageName: p.l.PackageName,
		FileName:    p.l.FileName,
	}

	if !p.peekTokenIs(token.LBRACE) {
		if p.peekTokenIs(token.EXTENDS) {
			p.nextToken()
			p.expectPeek(token.IDENT)
			classExpression.Parents = []*ast.Identifier{}
			classExpression.Parents = append(classExpression.Parents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			if p.peekTokenIs(token.COMMA) {
				p.nextToken()
				p.expectPeek(token.IDENT)
				classExpression.Parents = append(classExpression.Parents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			}
		}
		if p.peekTokenIs(token.IMPLEMENT) {
			p.nextToken()
			p.expectPeek(token.IDENT)
			classExpression.Interface = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}
	}
	p.nextToken()

	if !p.peekTokenIs(token.RBRACE) {
		letStatements := []*ast.LetStatement{}
		for p.peekTokenIs(token.LET) {
			p.nextToken()
			letStatement := p.parseLetStatement().(*ast.LetStatement)
			letStatements = append(letStatements, letStatement)
		}
		classExpression.LetStatements = letStatements

		functionStatemens := []*ast.FunctionLiteral{}
		for p.peekTokenIs(token.FUNCTION) {
			p.nextToken()
			functionStatement := p.parseFunctionLiteral().(*ast.FunctionLiteral)
			functionStatemens = append(functionStatemens, functionStatement)
			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}
		}
		classExpression.Functions = functionStatemens
	}
	p.expectPeek(token.RBRACE)
	return classExpression
}

func (p *Parser) parseInterfaceExpress() ast.Expression {
	interfaceExpress := &ast.InterfaceExpress{
		Token: p.curToken,
	}
	p.expectPeek(token.IDENT)
	interfaceExpress.Name = ast.Identifier{
		Token:       p.curToken,
		Value:       p.curToken.Literal,
		FileName:    p.l.FileName,
		PackageName: p.l.PackageName,
	}

	if p.peekTokenIs(token.EXTENDS) {
		p.nextToken()
		p.expectPeek(token.IDENT)
		interfaceExpress.Parents = []*ast.Identifier{}
		interfaceExpress.Parents = append(interfaceExpress.Parents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
			p.expectPeek(token.IDENT)
			interfaceExpress.Parents = append(interfaceExpress.Parents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		}
	}
	p.nextToken()

	if !p.peekTokenIs(token.RBRACE) {
		functionStatemens := []*ast.FunctionLiteral{}
		for p.peekTokenIs(token.FUNCTION) {
			p.nextToken()
			functionStatement := p.parseFunctionLiteral().(*ast.FunctionLiteral)
			functionStatemens = append(functionStatemens, functionStatement)
		}
		interfaceExpress.Functions = functionStatemens
	}
	p.expectPeek(token.RBRACE)
	return interfaceExpress
}

func (p *Parser) parseNewExpression() ast.Expression {
	newExpression := &ast.ObjectExpress{
		Token: p.curToken,
	}
	p.expectPeek(token.IDENT)
	newExpression.Class = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	p.expectPeek(token.LPAREN)
	newExpression.Parameters = p.parseExpressionList(token.RPAREN)
	return newExpression
}

func (p *Parser) parseDeferExpression() ast.Expression { // same as parseBlockStatement , except IsDeferBlock: true
	block := &ast.BlockStatement{Token: p.curToken, IsDeferBlock: true}
	block.Statements = []ast.Statement{}
	p.nextToken()
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseQuestionExpression(left ast.Expression) ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken, Condition: left}
	expression.Consequence = &ast.BlockStatement{
		Statements: []ast.Statement{},
	}
	p.nextToken()
	stmt := p.parseStatement()
	expression.Consequence.Statements = append(expression.Consequence.Statements, stmt)
	p.nextToken()
	if p.curTokenIs(token.COLON) {
		expression.Alternative = &ast.BlockStatement{
			Statements: []ast.Statement{},
		}
		p.nextToken()
		alterStmt := p.parseStatement()
		expression.Alternative.Statements = append(expression.Alternative.Statements, alterStmt)
	}
	return expression
}
