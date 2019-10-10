package parser

import (
	"fmt"
	
	"github.com/shavit/go-interpreter/ast"
	"github.com/shavit/go-interpreter/lexer"
	"github.com/shavit/go-interpreter/token"
)

//
// Parser
//
// The current and peek tokens acts exactly like in the lexer, except
//  they point to tokens instead of characters
type Parser struct {
	l *lexer.Lexer
	
	currentToken token.Token
	peekToken token.Token

	errors []string
}

// New creates a new parser from the lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},
	}

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
	default:
		return nil
	}
}

// parseLetStatement creates a let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	// Check the peek only after the statement was created
	//  since this function will change the parser state
	if !p.expectPeek(token.IDENT){
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	
	if !p.expectPeek(token.ASSIGN){
		return nil
	}

	// Skip the semicolon
	// This will be replaced later, once the 
	for !p.currentTokenIs(token.SEMICOLON){
		p.nextToken()
	}

	return stmt
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
