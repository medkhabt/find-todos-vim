package parser

import (
	"testing"

	"github.com/medkhabt/todoprs/token"
)

func TestFirstTerminalNode(t *testing.T) {
	input := &TerminalNode{&token.Token{token.CHAR, "a"}}
	validation := []*token.Token{&token.Token{token.CHAR, "a"}}

	parser := New()
	rslt := parser.first(input)
	for i, tt := range validation {
		tok := rslt[i]
		if rslt[i].Type != tt.Type {
			t.Fatalf("test[%d]: token type don't match. expected=%q, got=%q", i, tt.Type, tok.Type)
		}
		if rslt[i].Value != tt.Value {
			t.Fatalf("test[%d]: token value don't match. expected=%q, got=%q", i, tt.Value, tok.Value)
		}
	}

}
