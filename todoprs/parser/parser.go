package parser

import (
	"fmt"

	"github.com/medkhabt/todoprs/token"
)

type Parser struct {
}

type Node interface {
	productions() ([][]Node, error)
	isTerminal() bool
	getToken() (*token.Token, error)
}

type TerminalNode struct {
	token *token.Token
}

type NonTerminalNode struct {
	prods [][]Node
}

func New() *Parser {
	return &Parser{}
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

func (p *Parser) first(n Node) ([]*token.Token, error) {
	if n.isTerminal() {
		tk, err := n.getToken()
		if err != nil {
			return nil, err
		}
		return []*token.Token{tk}, nil
	} else {
		toks := []*token.Token{}
		productions, err := n.productions()
		if err != nil {
			return nil, err
		}
		for _, production := range productions {
			i := 0
			c := true
			for c && i < len(production) {
				c = false
				fy, err := p.first(production[i])
				if err != nil {
					return nil, err
				}
				if token.ExistEpsTk(fy) {
					c = true
					i += 1
				}
				for _, tok := range fy {
					if tok.Type != token.EPSILON {
						// TODO make it a set ? duplicate tokens
						toks = append(toks, tok)
					}
				}
				// Means that the last node of the production can also be reduced to epsilon so epsilon is also a first.
				if i == len(production) {
					toks = append(toks, &token.Token{token.EPSILON, ""})
				}
			}
		}
		return toks, nil
	}
}

/* func parse() error {
	// Choose node-production
	for i, n := range prod {
		if !n.isTerminal() {
			n.parse()
		} else if n == token.ILLEGAL { // change Illegal with appropriate token

		} else {
			return fmt.Errof("Parse Error in Node")
		}
	}
}
*/
