package ast

import (
	"bytes"

	"github.com/shavit/go-interpreter/token"
)

type Node interface {
	// TokenLiteral returns the literal value of a token
	//  which is a string.
	// It will be used for debugging and testing
	TokenLiteral() string

	// String prints AST nodes for debugging
	String() string
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

// String prints AST nodes for debugging
func (p *Program) String() string {
	var buf bytes.Buffer

	for _, s := range p.Statements {
		buf.WriteString(s.String())
	}

	return buf.String()
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

// String returns teh string value of the identifier
func (i *Identifier) String() string {
	return i.Value
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

// String prints AST nodes for debugging
func (l *LetStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(l.TokenLiteral() + " ")
	buf.WriteString(l.Name.String() + " = ")

	if l.Value != nil {
		buf.WriteString(l.Value.String())
	}
	buf.WriteString(";")

	return buf.String()
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

// String prints AST nodes for debugging
func (rs *ReturnStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		buf.WriteString(rs.ReturnValue.String())
	}

	buf.WriteString(";")

	return buf.String()
}

// ExpressionStatement allows the program to parse one line expressions
//  for example: `a + b` is an expression statement, and the result is
//  not assigned to a variable
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// statementNode is a helper that checks that this is a statement
func (es *ExpressionStatement) statementNode() {
}

// TokenLiteral returns the token literal
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String prints AST nodes for debugging
func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}

	return es.Expression.String()
}

// IntegerLiteral implements the Expression interface
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// expressionNode() returns the expression node
func (il *IntegerLiteral) expressionNode() {
}

// TokenLiteral() returns the integer token literal
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String() returns a string representation of Integer Literal
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
