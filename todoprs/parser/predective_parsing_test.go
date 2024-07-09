package parser

import (
	"testing"

	"github.com/medkhabt/todoprs/token"
)

func TestDepth2Grammar(t *testing.T) {
	// Grammar
	a := &TerminalNode{&token.Token{token.CHAR, "a"}}
	b := &TerminalNode{&token.Token{token.CHAR, "b"}}
	c := &TerminalNode{&token.Token{token.CHAR, "c"}}
	d := &TerminalNode{&token.Token{token.CHAR, "d"}}
	e := &TerminalNode{&token.Token{token.CHAR, "e"}}
	eps := &TerminalNode{&token.Token{token.EPSILON, ""}}

	B := (&NonTerminalNode{"B", [][]Node{}}).addProduction([]Node{b})
	D := (&NonTerminalNode{"D", [][]Node{}}).addProduction([]Node{eps})
	E := (&NonTerminalNode{"E", [][]Node{}}).addProduction([]Node{e})
	A := (&NonTerminalNode{"A", [][]Node{}}).addProduction([]Node{D, E}).addProduction([]Node{a})
	D.addProduction([]Node{d, A})
	S := (&NonTerminalNode{"S", [][]Node{}}).addProduction([]Node{A, B}).addProduction([]Node{c, D})

	p := New(S)
	parsingTable, err := p.makeParsingTable()
	if err != nil {
		t.Fatalf("This test is not your problem , you gotta fix your parsing table !: %s", err)
	}
	input := []*token.Token{
		&token.Token{token.CHAR, "a"},
		&token.Token{token.CHAR, "b"},
		&token.Token{token.EOF, ""},
	}
	err = p.PredictiveParsing(input, parsingTable)
	if err != nil {
		t.Fatalf("error from predective parsing %s.", err)
	}
}
