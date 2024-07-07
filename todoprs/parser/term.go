package parser

import (
	"fmt"

	"github.com/medkhabt/todoprs/token"
)

type TerminalNode struct {
	token *token.Token
}

func (tn *TerminalNode) productions() ([][]Node, error) {
	return nil, fmt.Errorf("TerminalNode has no productions.")
}
func (tn *TerminalNode) getToken() (*token.Token, error) {
	return tn.token, nil
}
func (tn *TerminalNode) isTerminal() bool {
	return true
}
