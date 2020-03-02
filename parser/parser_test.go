package parser

import (
	"fmt"
	"testing"

	"github.com/shavit/go-interpreter/ast"
	"github.com/shavit/go-interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 7;
let y = 11;
let someIdentifier = 390123;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Error("Found nil while expecting a Program")
		return
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Found %d statements, while expecting 3", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"someIdentifier"},
	}

	for i, item := range tests {
		stmt := program.Statements[i]
		if !testLetStatements(t, stmt, item.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Found %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("Found %q, while expecting `let`", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("Found %T, while expecting `*ast.LetStatement", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("Found %s, while expecting %s", letStmt.Name.Value, name)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Found %s, while expecting %s", letStmt.Name, name)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 7;
return 11;
return 31234;
return -1;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("Found %d, while expecting 4 program statements", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Invalid statement. Found %q, while expecting returnStmt.TokenLiteral to be \"return\"", returnStmt.TokenLiteral())
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Incorrect statement. Found %q, while expecting returnStmt.TokenLiteral to be \"return\"", returnStmt.TokenLiteral())
		}
	}
}

// Identifiers will evaluate their value. They produce a value like any
//  other expressions
func TestIndentifierExpression(t *testing.T) {
	input := "someIdentifier;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Found %d, while expecting 1 statement", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Found %T, while expecting ast.ExpressionStatement in program.Statements[0]", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Found %T, while expecting exp to be *ast.Identifier", stmt.Expression)
	}

	if ident.Value != "someIdentifier" {
		t.Errorf("Found %s, while expecting ident.Value to be  %s", ident.Value, "someIdentifier")
	}

	if ident.TokenLiteral() != "someIdentifier" {
		t.Errorf("Found %s, while expecting ident.TokenLiteral to be %s", ident.TokenLiteral(), "someIdentifier")
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "7;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Go %d, while expecting 1 statement", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Got %T, while expecting program.Statements[0] to be ast.ExpressionStatement", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Got %T, while expecting stmt.Expression to be *ast.IntegerLiteral", stmt.Expression)
	}
	if literal.Value != 7 {
		t.Errorf("Got %d, while expecting literal.Value to be %d", literal.Value, 7)
	}
	if literal.TokenLiteral() != "7" {
		t.Errorf("Got %s, while expecting literal.TokenLiteral to be %s", literal.TokenLiteral(), "7")
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!4", "!", 4},
		{"-21", "-", 21},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Got %d, while expecting program.Statements to have %d statements for input %s", len(program.Statements), 1, tt.input)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Got %T, while expecting program.Statements[0] to be ast.ExpressionStatement", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Got %T, while expecting stmt to be ast.PrefixExpression", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Got %s, while expecting exp.Operator to be %s", exp.Operator, tt.operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) (ok bool) {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Got %T, while expecting *ast.IntegerLiteral", il)
		return
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Got %s, while expecting %d", integ.TokenLiteral(), value)
		ok = false
		return
	}

	ok = true

	return ok
}

func TsetParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"4 + 9;", 4, "+", 9},
		{"5 - 6;", 5, "-", 6},
		{"6 * 3;", 6, "*", 3},
		{"10 / 2;", 10, "/", 2},
		{"7 > 8;", 7, ">", 8},
		{"9 < 10;", 9, "<", 10},
		{"2 == 7;", 2, "==", 7},
		{"4 != 4;", 4, "!=", 4},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Got %d, while expecting program.Statements to contain %d statements", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Got %T, while expecting program.Statements[0] to be ast.ExpressionStatement", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Got %T, while expecting exp to be ast.InfixExpression", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Got %s, while expecting exp.Operator to be %s", exp.Operator, tt.operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOPeratorPrecendenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Found %q, while expecting %q", actual, tt.expected)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", "true"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Found %q, while expecting %q", actual, tt.expected)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Found %T, whie expecting exp to be *ast.Identifier", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("Found %s, while expecting ident.Value to be %s", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("Found %s, while expecting ident.TokenLiteral to be %s", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expect interface{}) bool {
	switch v := expect.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("Unmatched type %T. Test not implemented.", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Found %T(%s), while expecting exp to be ast.OperatorExpression", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("Found %q, while expecting exp.Operator to be %s", opExp.Operator, opExp)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
