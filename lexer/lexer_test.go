package lexer

import (
	"testing"

	"github.com/shavit/go-interpreter/token"
)

func TestGetNextToken(t *testing.T) {
	input := `let twelve = 12;
let four = 4;

let add = fn(x, y){
  x + y;
};

let result = add(twelve,four);
!-/*9;
3 < 10 > 7;

if (0<42){
  return true;
} else {
return false;
}

18 == 18;
19 != 17;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "twelve"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "four"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "twelve"},
		{token.COMMA, ","},
		{token.IDENT, "four"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.INT, "3"},
		{token.LT, "<"},
		{token.INT, "1"},
		{token.INT, "0"},
		{token.GT, ">"},
		{token.INT, "7"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "0"},
		{token.LT, "<"},
		{token.INT, "4"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "1"},
		{token.INT, "8"},
		{token.EQ, "=="},
		{token.INT, "1"},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.INT, "9"},
		{token.NOT_EQ, "!="},
		{token.INT, "1"},
		{token.INT, "7"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lxr := New(input)

	for i, item := range tests {
		tkn := lxr.NextToken()

		if tkn.Type != item.expectedType {
			t.Fatalf("Error at %d: Got: %q, while epxecting: %q", i, item.expectedType, tkn.Type)
		}

		if tkn.Literal != item.expectedLiteral {
			t.Fatalf("Error at %d: Got: %q, while epxecting: %q", i, item.expectedLiteral, tkn.Literal)
		}
	}
}
