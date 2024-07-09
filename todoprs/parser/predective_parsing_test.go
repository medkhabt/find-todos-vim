package parser

import (
	"testing"

	"github.com/medkhabt/todoprs/token"
)

func TestDepth2Grammar(t *testing.T) {
	// Grammar
	a := &TerminalNode{&token.Token{token.CHAR, "a"}}
	b := &TerminalNode{&token.Token{token.SLASH, ""}}
	c := &TerminalNode{&token.Token{token.COLON, ""}}
	d := &TerminalNode{&token.Token{token.DOT, ""}}
	e := &TerminalNode{&token.Token{token.TODO, ""}}
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
		&token.Token{token.SLASH, ""},
		&token.Token{token.EOF, ""},
	}
	err = p.PredictiveParsing(input, parsingTable)
	if err != nil {
		t.Fatalf("error from predective parsing %s.", err)
	}
}

func TestOnRipGrepGrammar(t *testing.T) {
	c := &TerminalNode{&token.Token{token.CHAR, ""}}
	d := &TerminalNode{&token.Token{token.DIGIT, ""}}
	colon := &TerminalNode{&token.Token{token.COLON, ""}}
	slash := &TerminalNode{&token.Token{token.SLASH, ""}}
	dot := &TerminalNode{&token.Token{token.DOT, ""}}
	todo := &TerminalNode{&token.Token{token.TODO, ""}}
	space := &TerminalNode{&token.Token{token.SPACE, ""}}
	eps := &TerminalNode{&token.Token{token.EPSILON, ""}}
	newline := &TerminalNode{&token.Token{token.NEWLINE, ""}}

	C0 := (&NonTerminalNode{"C0", [][]Node{}}).addProduction([]Node{c}).addProduction([]Node{d})
	// Maaaan this is not pretty.. why dot at the end golang :/
	C := (&NonTerminalNode{"C", [][]Node{}}).
		addProduction([]Node{C0}).
		addProduction([]Node{dot}).
		addProduction([]Node{colon}).
		addProduction([]Node{space}).
		addProduction([]Node{slash})
	T1 := (&NonTerminalNode{"T1", [][]Node{}})
	T1.addProduction([]Node{C, T1}).addProduction([]Node{eps})
	T := (&NonTerminalNode{"T", [][]Node{}}).addProduction([]Node{C, T1})

	L1 := (&NonTerminalNode{"L1", [][]Node{}})
	L1.addProduction([]Node{d, L1}).addProduction([]Node{eps})
	L := (&NonTerminalNode{"L", [][]Node{}}).addProduction([]Node{colon, d, L1, colon})

	E := (&NonTerminalNode{"E", [][]Node{}}).addProduction([]Node{todo})
	E.addProduction([]Node{C, E})

	P1 := (&NonTerminalNode{"P1", [][]Node{}})
	P0 := (&NonTerminalNode{"P0", [][]Node{}}).addProduction([]Node{C0, P1})
	P1.addProduction([]Node{C0, P1}).
		addProduction([]Node{slash, C0, P1}).
		addProduction([]Node{dot, C0, P1}).
		addProduction([]Node{eps})
	St := (&NonTerminalNode{"St", [][]Node{}}).addProduction([]Node{P0, L, slash, slash, E, T})
	S1 := (&NonTerminalNode{"S1", [][]Node{}}).addProduction([]Node{newline, St}).addProduction([]Node{eps})
	S := (&NonTerminalNode{"S", [][]Node{}}).addProduction([]Node{St, S1})

	p := New(S)
	parsingTable, err := p.makeParsingTable()
	if err != nil {
		t.Fatalf("This test is not your problem , you gotta fix your parsing table !: %s", err)
	}
	// t/t.t:33://TODO  t  t\ntt.t:33://TODO  t  t$
	input := []*token.Token{
		&token.Token{token.CHAR, "t"},
		&token.Token{token.SLASH, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.DOT, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.COLON, ""},
		&token.Token{token.DIGIT, "3"},
		&token.Token{token.DIGIT, "3"},
		&token.Token{token.COLON, ""},
		&token.Token{token.SLASH, ""},
		&token.Token{token.SLASH, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.TODO, ""},
		&token.Token{token.SPACE, ""},
		&token.Token{token.SPACE, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.SPACE, ""},
		&token.Token{token.SPACE, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.NEWLINE, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.DOT, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.COLON, ""},
		&token.Token{token.DIGIT, "3"},
		&token.Token{token.DIGIT, "3"},
		&token.Token{token.COLON, ""},
		&token.Token{token.SLASH, ""},
		&token.Token{token.SLASH, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.TODO, ""},
		&token.Token{token.SPACE, ""},
		&token.Token{token.SPACE, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.SPACE, ""},
		&token.Token{token.SPACE, ""},
		&token.Token{token.CHAR, "t"},
		&token.Token{token.EOF, ""},
	}
	err = p, PredictiveParsing(input, parsingTable)
	if err != nil {
		t.Fatalf("error from predective parsing %s.", err)
	}
}
