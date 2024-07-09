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

// a is colon token, b is DIGIT, c SLASH, d TODO, e DOT
// S -> AB , S -> cD , A -> DE, A -> a , D -> dA, D -> \eps , E -> e , B -> b,

// __ | a       | b          | c        | d           | e         | $         |
// ---|---------|------------|----------|-------------|-----------|-----------|
// A  | A -> :  |            |          | A -> D E    | A -> D E  |           |
// ---|---------|------------|----------|-------------|-----------|-----------|
// B  |         | B -> DIGIT |          |             |           |           |
// ---|---------|------------|----------|-------------|-----------|-----------|
// D  |         |            |          | D -> TODO A | D -> \eps |           |
// ---|---------|------------|----------|-------------|-----------|-----------|
// E  |         |            |          |             | E -> .    |           |
// ---|-------- |----------- |----------|-------------|-----------|-----------|
// S  | S-> A B |            | S -> / D | S -> A B    | S -> A B  |           |
// ---------------------------------------------------------------------------|
func TestWithDepth2(t *testing.T) {
	a := &TerminalNode{&token.Token{token.COLON, ""}}
	b := &TerminalNode{&token.Token{token.DIGIT, ""}}
	c := &TerminalNode{&token.Token{token.SLASH, ""}}
	d := &TerminalNode{&token.Token{token.TODO, ""}}
	e := &TerminalNode{&token.Token{token.DOT, ""}}
	eps := &TerminalNode{&token.Token{token.EPSILON, ""}}

	B := (&NonTerminalNode{"B", [][]Node{}}).addProduction([]Node{b})
	D := (&NonTerminalNode{"D", [][]Node{}}).addProduction([]Node{eps})
	E := (&NonTerminalNode{"E", [][]Node{}}).addProduction([]Node{e})
	A := (&NonTerminalNode{"A", [][]Node{}}).addProduction([]Node{D, E}).addProduction([]Node{a})
	D.addProduction([]Node{d, A})
	S := (&NonTerminalNode{"S", [][]Node{}}).addProduction([]Node{A, B}).addProduction([]Node{c, D})

	validation := map[Transition][]Node{
		Transition{A.name, a.token.Type}: []Node{a},
		Transition{A.name, d.token.Type}: []Node{D, E},
		Transition{A.name, e.token.Type}: []Node{D, E},
		Transition{B.name, b.token.Type}: []Node{b},
		Transition{D.name, d.token.Type}: []Node{d, A},
		Transition{D.name, e.token.Type}: []Node{eps},
		Transition{E.name, e.token.Type}: []Node{e},
		Transition{S.name, c.token.Type}: []Node{c, D},
		Transition{S.name, d.token.Type}: []Node{A, B},
		Transition{S.name, e.token.Type}: []Node{A, B},
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
			t.Fatalf("[%s] Transition[%s,%+v]: Prod didn't match. expected=%v, got=%v", testName, tr.nodeName, tr.transitorType, ne, nr)
		}
	}
}
