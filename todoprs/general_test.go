package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/medkhabt/todoprs/parser"
)

func TestGeneral(t *testing.T) {
	S := parser.RgGrammar()
	p := parser.New(S)
	input := "/t /t.t  :30://TODO Add feature one.  \n/t.ttt:324://TODO Fix this bug."
	err := p.Parse(input)
	if err != nil {
		t.Fatalf("error from predective parsing %s.", err)
	}
}
func TestGeneralReal(t *testing.T) {
	path := os.Getenv("TODO_PATH_PRJ")
	cmd := exec.Command("rg", "-nw", "TODO", path)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("problem with the output from the command")
	}
	S := parser.RgGrammar()
	p := parser.New(S)
	err = p.Parse(string(out))
	if err != nil {
		t.Fatalf("error from predective parsing %s.", err)
	}
}
