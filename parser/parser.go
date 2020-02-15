package parser

import (
	"fmt"

	"github.com/shavit/go-interpreter/ast"
	"github.com/shavit/go-interpreter/lexer"
	"github.com/shavit/go-interpreter/token"
)

// iota because the order of the expressions matter
//  for example, + have lower precedence than *
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // someFunc(x)
)

//
// Parser
//
// The current and peek tokens acts exactly like in the lexer, except
//  they point to tokens instead of characters
type Parser struct {
	l *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New creates a new parser from the lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// Two tokens to set the current and peek tokens
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns the parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError check for errors in the next token
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf(`Found %s, whlie expecting the next token to be %s`, p.peekToken.Type, t)
	p.errors = append(p.errors, msg)
}

// nextToken advance the parser to the next token
func (p *Parser) nextToken() {
	// Take the peek token from this parser
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram creates a tree of statements from the lexer
func (p *Parser) ParseProgram() *ast.Program {
	program := new(ast.Program)
	program.Statements = []ast.Statement{}

	// Iterate through the tokens and create statements
	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement is a helper to parse a statement
// It returns the Statement interface
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement creates a let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	// Check the peek only after the statement was created
	//  since this function will change the parser state
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// Skip the semicolon
	// This will be replaced later, once the
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement creates a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	// Go to the next token after the return statement
	p.nextToken()

	// Parse until the semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses expressions
// This is the default parser
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	// Pass the lowest since nothing was parsed yet
	stmt.Expression = p.parseExpression(LOWEST)

	// The semicolons are optional, to make the REPL simpler
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpression parse individual expression from a statement
// it checks if there is a parsing function for the current token
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

// currentTokenIs check the current token against a type
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs check the peek token against a type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek advance to the next token if the peek match a type
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

type prefixParseFn func() ast.Expression

// registerPrefix adds a parse function to the prefix function map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

type infixParseFn func(ast.Expression) ast.Expression

// registerInfix adds a parse function to the infix function map
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parseIdentifier returns an identifier with the current token, and
//  the literal token value
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}
