package parser

import "github.com/medkhabt/todoprs/token"

type Parser struct {
}

type Node interface {
	productions() [][]Node
	isTerminal() bool
	getToken() *token.Token
}

type TerminalNode struct {
	token *token.Token
}

func New() *Parser {
	return &Parser{}
}
func (tn *TerminalNode) productions() [][]Node {
	return nil
}
func (tn *TerminalNode) getToken() *token.Token {
	return tn.token
}
func (tn *TerminalNode) isTerminal() bool {
	return true
}

func (p *Parser) first(n Node) []*token.Token {
	if n.isTerminal() {
		return []*token.Token{n.getToken()}
	} else {
		toks := []*token.Token{}
		for _, production := range n.productions() {
			i := 0
			c := true
			for c && i < len(production) {
				c = false
				fy := p.first(production[i])
				if token.ExistEpsTk(fy) {
					c = true
					i += 1
				}
				if len(production) == 1 && len(fy) == 1 && fy[0].Type == token.EPSILON {
					// A -> \epsilon, add the epsilon as a first
					toks = append(toks, fy[0])
				}
				for _, tok := range fy {
					if tok.Type != token.EPSILON {
						// TODO make it a set ? duplicate tokens
						toks = append(toks, tok)
					}
				}
			}
		}
		return toks
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
