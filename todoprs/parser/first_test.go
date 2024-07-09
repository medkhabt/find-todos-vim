package parser

import (
	"testing"

	"github.com/medkhabt/todoprs/token"
)

func TestFirstTerminalNode(t *testing.T) {
	input := &TerminalNode{&token.Token{token.CHAR, "a"}}
	validation := []*token.Token{&token.Token{token.CHAR, "a"}}

	parser := New(input)
	rslt, err := parser.first(input)
	if err != nil {
		t.Fatalf("result returned the following error : %s", err)
	}
	for i, tt := range validation {
		tok := rslt[i]
		if rslt[i].Type != tt.Type {
			t.Fatalf("token[%d]: token type don't match. expected=%q, got=%q", i, tt.Type, tok.Type)
		}
	}

}

// TODO order of tokens should when checking the result of first function
func TestFirstNonTerminalNode(t *testing.T) {
	tn := &TerminalNode{&token.Token{token.CHAR, "a"}}
	tnb := &TerminalNode{&token.Token{token.SLASH, ""}}
	tnc := &TerminalNode{&token.Token{token.COLON, ""}}
	tnd := &TerminalNode{&token.Token{token.DOT, ""}}
	eps := &TerminalNode{&token.Token{token.EPSILON, ""}}
	// A -> a    :::::> First(A) = {a}
	A0 := (&NonTerminalNode{}).addProduction([]Node{tn})
	// A -> a | \epsilon  :::::> First(A) = {a, \epsilon}
	A1 := (&NonTerminalNode{}).addProduction([]Node{tn}).addProduction([]Node{eps})
	// A -> aB  :::::> First(A) = {a}
	B := (&NonTerminalNode{}).addProduction([]Node{tnb})
	A2 := (&NonTerminalNode{}).addProduction([]Node{tn, B})
	// A -> BC  ::: B -> b  :::::> First(A) = {b}
	B = (&NonTerminalNode{}).addProduction([]Node{tnb})
	C := (&NonTerminalNode{}).addProduction([]Node{tnc})
	A3 := (&NonTerminalNode{}).addProduction([]Node{B, C})
	// A -> BC  ::: B -> b| \epsilon ::: C -> c :::::> First(A) = {b,c}
	B = (&NonTerminalNode{}).addProduction([]Node{tnb}).addProduction([]Node{eps})
	C = (&NonTerminalNode{}).addProduction([]Node{tnc})
	A4 := (&NonTerminalNode{}).addProduction([]Node{B, C})
	// A -> BC | \epsilon ::: B -> b| \epsilon ::: C -> c :::::> First(A) = {b,c,\epsilon}
	B = (&NonTerminalNode{}).addProduction([]Node{tnb}).addProduction([]Node{eps})
	C = (&NonTerminalNode{}).addProduction([]Node{tnc})
	A5 := (&NonTerminalNode{}).addProduction([]Node{B, C}).addProduction([]Node{eps})
	// A -> BC ::: B -> b| \epsilon ::: C -> c | \epsilon :::::> First(A) = {b,c,\epsilon}
	B = (&NonTerminalNode{}).addProduction([]Node{tnb}).addProduction([]Node{eps})
	C = (&NonTerminalNode{}).addProduction([]Node{tnc}).addProduction([]Node{eps})
	A6 := (&NonTerminalNode{}).addProduction([]Node{B, C})
	// A -> BC ::: B -> b| \epsilon ::: C -> c | \epsilon
	// D -> d , Z -> AD :::::> First(Z) = {b,c,\epsilon, d}
	B = (&NonTerminalNode{}).addProduction([]Node{tnb}).addProduction([]Node{eps})
	C = (&NonTerminalNode{}).addProduction([]Node{tnc}).addProduction([]Node{eps})
	D := (&NonTerminalNode{}).addProduction([]Node{tnd})
	A7 := (&NonTerminalNode{}).addProduction([]Node{B, C})
	Z := (&NonTerminalNode{}).addProduction([]Node{A7, D})

	input := []Node{
		A0,
		A1,
		A2,
		A3,
		A4,
		A5,
		A6,
		Z,
	}
	validation := [][]*token.Token{
		[]*token.Token{&token.Token{token.CHAR, "a"}},
		[]*token.Token{&token.Token{token.CHAR, "a"}, &token.Token{token.EPSILON, ""}},
		[]*token.Token{&token.Token{token.CHAR, "a"}},
		[]*token.Token{&token.Token{token.SLASH, "b"}},
		[]*token.Token{&token.Token{token.SLASH, ""}, &token.Token{token.COLON, ""}},
		[]*token.Token{&token.Token{token.SLASH, ""}, &token.Token{token.COLON, ""}, &token.Token{token.EPSILON, ""}},
		[]*token.Token{&token.Token{token.SLASH, ""}, &token.Token{token.COLON, ""}, &token.Token{token.EPSILON, ""}},
		[]*token.Token{&token.Token{token.SLASH, ""}, &token.Token{token.COLON, ""}, &token.Token{token.DOT, ""}},
	}
	for k := 0; k < len(input); k += 1 {
		parser := New(input[k])
		rslt, err := parser.first(input[k])
		if err != nil {
			t.Fatalf("result returned the following error : %s", err)
		}
		for i, tt := range validation[k] {
			if len(rslt) != len(validation[k]) {
				t.Fatalf("result length %d diff from validation length %d for test %d.", len(rslt), len(validation[k]), k)
			}
			tok := rslt[i]
			if rslt[i].Type != tt.Type {
				t.Fatalf("test[%d] token[%d]: token type don't match. expected=%q, got=%q", k, i, tt.Type, tok.Type)
			}
		}
	}

}
