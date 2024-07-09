package parser

import (
	"runtime"
	"strings"
	"testing"

	"github.com/medkhabt/todoprs/token"
)

func TestFollowStartSymbol(t *testing.T) {
	// A shouldn't be relevant in this case
	a := &TerminalNode{&token.Token{token.CHAR, "a"}}
	A := (&NonTerminalNode{name: "A"}).addProduction([]Node{a})
	input := (&NonTerminalNode{name: "S"}).addProduction([]Node{a, A})
	// Should it be end of file or should i differenciate the $ from eof.
	validation := []*token.Token{&token.Token{token.EOF, ""}}

	parser := New(input)
	rslt, err := parser.follow(input)
	if err != nil {
		t.Fatalf("result returned the following error : %s", err)
	}
	assertEq(t, rslt, validation, getFunctionName())
}

func TestFollowInMiddle(t *testing.T) {
	// testing the rule A => \alpha B \beta , it should add first(\beta) to follow(B)
	a := &TerminalNode{&token.Token{token.CHAR, "a"}}
	b := &TerminalNode{&token.Token{token.CHAR, "b"}}
	c := &TerminalNode{&token.Token{token.CHAR, "c"}}
	B := (&NonTerminalNode{name: "B"}).addProduction([]Node{b})
	C := (&NonTerminalNode{name: "C"}).addProduction([]Node{c})
	input := (&NonTerminalNode{name: "S"}).addProduction([]Node{a, B, C})
	ct, err := c.getToken()
	if err != nil {
		t.Fatalf("Error when getting the token to setup the validation value :[%s].", err)
	}
	validation := []*token.Token{ct}

	parser := New(input)
	rslt, err := parser.follow(B)
	if err != nil {
		t.Fatalf("result returned the following error : %s", err)
	}
	assertEq(t, rslt, validation, getFunctionName())
}

func TestFollowWithBetaDerivableToEpsilon(t *testing.T) {
	// testing the rule A => \alpha B \beta , it should add first(\beta) to follow(B)
	a := &TerminalNode{&token.Token{token.CHAR, "a"}}
	b := &TerminalNode{&token.Token{token.SLASH, ""}}
	c := &TerminalNode{&token.Token{token.COLON, ""}}
	eps := &TerminalNode{&token.Token{token.EPSILON, ""}}
	B := (&NonTerminalNode{name: "B"}).addProduction([]Node{b})
	C := (&NonTerminalNode{name: "C"}).addProduction([]Node{c}).addProduction([]Node{eps})
	input := (&NonTerminalNode{name: "S"}).addProduction([]Node{a, B, C})
	ct, err := c.getToken()
	if err != nil {
		t.Fatalf("Error when getting the token to setup the validation value :[%s].", err)
	}
	validation := []*token.Token{ct, &token.Token{token.EOF, ""}}

	parser := New(input)
	rslt, err := parser.follow(B)
	if err != nil {
		t.Fatalf("result returned the following error : %s", err)
	}
	assertEq(t, rslt, validation, getFunctionName())
}

func TestFollowWithAlphaB(t *testing.T) {
	// testing the rule A => \alpha B \beta , it should add first(\beta) to follow(B)
	a := &TerminalNode{&token.Token{token.CHAR, "a"}}
	b := &TerminalNode{&token.Token{token.SLASH, ""}}
	B := (&NonTerminalNode{name: "B"}).addProduction([]Node{b})
	input := (&NonTerminalNode{name: "S"}).addProduction([]Node{a, B})
	validation := []*token.Token{&token.Token{token.EOF, ""}}

	parser := New(input)
	rslt, err := parser.follow(B)
	if err != nil {
		t.Fatalf("result returned the following error : %s", err)
	}
	assertEq(t, rslt, validation, getFunctionName())
}

func TestFollowWithDepth2(t *testing.T) {
	// testing the rule A => \alpha B \beta , it should add first(\beta) to follow(B)
	a := &TerminalNode{&token.Token{token.CHAR, ""}}
	b := &TerminalNode{&token.Token{token.SLASH, ""}}
	c := &TerminalNode{&token.Token{token.COLON, ""}}
	d := &TerminalNode{&token.Token{token.DOT, ""}}
	eps := &TerminalNode{&token.Token{token.EPSILON, ""}}
	B := (&NonTerminalNode{name: "B"}).addProduction([]Node{b})
	C := (&NonTerminalNode{name: "C"}).addProduction([]Node{c, B}).addProduction([]Node{eps})
	input := (&NonTerminalNode{name: "S"}).addProduction([]Node{a, B, C, d})
	ct, err := c.getToken()
	dt, err := d.getToken()
	if err != nil {
		t.Fatalf("Error when getting the token to setup the validation value :[%s].", err)
	}

	validation := []*token.Token{ct, dt}

	parser := New(input)
	rslt, err := parser.follow(B)
	if err != nil {
		t.Fatalf("result returned the following error : %s", err)
	}
	assertEq(t, rslt, validation, getFunctionName())
}
func assertEq(t *testing.T, rslt []*token.Token, validation []*token.Token, testName string) {
	for i, tt := range validation {
		tok := rslt[i]
		if rslt[i].Type != tt.Type {
			t.Fatalf("[%s] token[%d]: token type don't match. expected=%q, got=%q", testName, i, tt.Type, tok.Type)
		}
	}
}

func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return strings.Split(runtime.FuncForPC(pc).Name(), ".")[2]
	//return strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
}
