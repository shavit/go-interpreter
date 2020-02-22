package lexer

import (
	"github.com/shavit/go-interpreter/token"
)

type Lexer struct {
	input        string
	position     int  // Current character position (ch)
	readPosition int  // After current character
	ch           byte // Current character (in position)
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()

	return l
}

// readChar reads the next character
func (l *Lexer) readChar() {
	// Advance by 1 or assign NUL at the end
	if l.readPosition >= len(l.input) {
		l.ch = 0x0
	} else {
		l.ch = l.input[l.readPosition]
	}

	// Move to the next byte
	l.position = l.readPosition
	l.readPosition += 1
}

// peekChar reads a character without incrementing the position
// It is being use to peek the next character
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// NextToken gets the next token
func (l *Lexer) NextToken() token.Token {
	var tkn token.Token
	l.ignoreWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tkn = newToken(token.PLUS, l.ch)
	case '-':
		tkn = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(token.BANG, l.ch)
		}
	case '/':
		tkn = newToken(token.SLASH, l.ch)
	case '*':
		tkn = newToken(token.ASTERISK, l.ch)
	case '<':
		tkn = newToken(token.LT, l.ch)
	case '>':
		tkn = newToken(token.GT, l.ch)
	case ';':
		tkn = newToken(token.SEMICOLON, l.ch)
	case '(':
		tkn = newToken(token.LPAREN, l.ch)
	case ')':
		tkn = newToken(token.RPAREN, l.ch)
	case ',':
		tkn = newToken(token.COMMA, l.ch)
	case '{':
		tkn = newToken(token.LBRACE, l.ch)
	case '}':
		tkn = newToken(token.RBRACE, l.ch)
	case 0x0:
		tkn.Literal = ""
		tkn.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tkn.Literal = l.readIdentifier()
			tkn.Type = token.LookupIdent(tkn.Literal)
			return tkn
		} else if isDigit(l.ch) {
			tkn.Type = token.INT
			tkn.Literal = l.readNumber()
			return tkn
		} else {
			tkn = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tkn
}

// newToken creates a new token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// Ignore white space
func (l *Lexer) ignoreWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// isLetter checks if the current byte is a letter
// it checks if the byte in the range of [a-zA-Z_]
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// readNumber reads a number
// It only supports integers
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// ifDigit checks if the current byte is a digit
// it checks if the byte in the range of [0-9]
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
