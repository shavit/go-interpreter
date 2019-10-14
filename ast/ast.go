package ast

import (
	"github.com/shavit/go-interpreter/token"
)

type Node interface {
	// TokenLiteral returns the literal value of a token
	//  which is a string.
	// It will be used for debugging and testing
	TokenLiteral() string
}

type Statement interface {
	Node
	// statementNode checks that this is a statement node
	statementNode()
}

type Expression interface {
	Node
	// expressionNode checks that this is an expression node
	expressionNode()
}

type Program struct {
	// Slice of AST Statement nodes
	Statements []Statement
}

// TokenLiteral returns the literal value of a token
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token
	Value string
}

// expressionsNode is a helper that checks that this is an expression
func (li *Identifier) expressionNode() {
}

// TokenLiteral returns the token value
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatement struct {
	// The LET token
	Token token.Token
	Name  *Identifier
	Value Expression
}

// statementNode is a helper that checks that this is a statement
func (l *LetStatement) statementNode() {
}

// TokenLiteral returns the token value
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

// statementNode is a helper that checks that this is a statement
func (rs *ReturnStatement) statementNode() {
}

// TokenLiteral returns the token literal
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
