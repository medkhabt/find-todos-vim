package parser

import (
	"fmt"

	"github.com/medkhabt/todoprs/token"
)

type NonTerminalNode struct {
	name  string
	prods [][]Node
}

func (ntn *NonTerminalNode) productions() ([][]Node, error) {
	return ntn.prods, nil
}
func (ntn *NonTerminalNode) getToken() (*token.Token, error) {
	return nil, fmt.Errorf("getToken() not available for a NonTerminalNode.")
}
func (ntn *NonTerminalNode) isTerminal() bool {
	return false
}
func (ntn *NonTerminalNode) addProduction(production []Node) *NonTerminalNode {
	ntn.prods = append(ntn.prods, production)
	return ntn
}
