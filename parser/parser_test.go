package parser

import (
	"testing"

	"github.com/shavit/go-interpreter/ast"
	"github.com/shavit/go-interpreter/lexer"
)

func TestLetStatements(t *testing.T){
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

func checkParserErrors(t *testing.T, p *Parser){
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
