package parser

import (
	"testing"

	"github.com/medkhabt/todoprs/comparator"
	"github.com/medkhabt/todoprs/token"
)

// for each A -> \alpha
//
//	terminal a in FIRST(\alph), add A -> \alph to M[A,a]
//	if eps in FIRST(\alph), then foreach b terminal in FOLLOW(A), add A -> \alpha to M[A, b]

// S -> AB , S -> cD , A -> DE, A -> a , D -> dA, D -> \eps , E -> e , B -> b,

// __ | a      | b         | c       | d       | e         | $         |
// ---|--------|-----------|---------|---------|-----------|-----------|
// A  | A -> a |           |         | A -> DE | A -> DE   |           |
// ---|--------|-----------|---------|---------|-----------|-----------|
// B  |        | B -> b    |         |         |           |           |
// ---|--------|-----------|---------|---------|-----------|-----------|
// D  |        |           |         | D -> dA | D -> \eps |           |
// ---|--------|-----------|---------|---------|-----------|-----------|
// E  |        |           |         |         | E -> e    |           |
// ---|--------|-----------|---------|---------|-----------|-----------|
// S  | S-> AB |           | S -> cD | S -> AB | S -> AB   |           |
// --------------------------------------------------------------------|
func TestWithDepth2(t *testing.T) {
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

	validation := map[Transition][]Node{
		Transition{A.name, a.token.Value}: []Node{a},
		Transition{A.name, d.token.Value}: []Node{D, E},
		Transition{A.name, e.token.Value}: []Node{D, E},
		Transition{B.name, b.token.Value}: []Node{b},
		Transition{D.name, d.token.Value}: []Node{d, A},
		Transition{D.name, e.token.Value}: []Node{eps},
		Transition{E.name, e.token.Value}: []Node{e},
		Transition{S.name, c.token.Value}: []Node{c, D},
		Transition{S.name, d.token.Value}: []Node{A, B},
		Transition{S.name, e.token.Value}: []Node{A, B},
	}

	parser := New(S)
	rslt, err := parser.makeParsingTable()
	if err != nil {
		t.Fatalf("result returned the following error : %s", err)
	}
	assertEqMap(t, rslt, validation, getFunctionName())
}

// Make this generic  ?
func assertEqMap(t *testing.T, rslt map[Transition][]Node, validation map[Transition][]Node, testName string) {
	for tr, ne := range validation {
		nr := rslt[tr]
		if !comparator.CmpSlice(nr, ne) {
			t.Fatalf("[%s] Transition[%s,%+v]: Prod didn't match. expected=%v, got=%v", testName, tr.nodeName, tr.transitorValue, ne, nr)
		}
	}
}
