package ast

import (
	"testing"

	"github.com/shavit/go-interpreter/token"
)

func TestGetProgramString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "socketAddr"},
					Value: "socketAddr",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "addr"},
					Value: "addr",
				},
			},
		},
	}

	if program.String() != "let socketAddr = addr;" {
		t.Errorf("Wrong string value. Got %q", program.String())
	}
}
