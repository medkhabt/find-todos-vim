package parser

import (
	"fmt"

	"github.com/medkhabt/todoprs/token"
)

type Parser struct {
	startNode Node
}

type Node interface {
	productions() ([][]Node, error)
	isTerminal() bool
	getToken() (*token.Token, error)
}

// TODO this is slice based queue, if i have time or it is need or i just want to, than i should change this to ring-buffer based queue
type Queue[E any] struct {
	list  []E
	front int
}

func (q *Queue[E]) enqueue(e E) {
	q.list = append(q.list, e)
}
func (q *Queue[E]) dequeue() (E, error) {
	if q.list == nil || q.front >= len(q.list) {
		var e E
		return e, fmt.Errorf("Empty queue")
	} else {
		q.front += 1
		return q.list[q.front-1], nil
	}
}

func (q *Queue[E]) empty() bool {
	return q.front >= len(q.list)
}
func New(startNode Node) *Parser {
	return &Parser{startNode}
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

func (p *Parser) isStart(n Node) bool {
	return n == p.startNode
}

// A -> \alpha B
func (p *Parser) alphaB(result []*token.Token, A Node) ([]*token.Token, error) {
	followA, err := p.follow(A)
	if err != nil {
		return nil, err
	}
	for _, lf := range followA {
		result = append(result, lf)
	}
	return result, nil
}

func (p *Parser) alphaBbeta(result []*token.Token, currNode Node, prd []Node, i int) ([]*token.Token, error) {
	// alpha B beta
	c := true
	lenPrd := len(prd)
	// targeting the beta
	i += 1
	for c {
		c = false
		if i >= lenPrd {
			var err error
			result, err = p.alphaB(result, currNode)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		firsts, err := p.first(prd[i])
		if err != nil {
			return nil, err
		}
		for _, f := range firsts {
			if f.Type == token.EPSILON {
				c = true
			} else {
				result = append(result, f)
			}
		}
		i += 1
	}
	return result, nil
}

// BFS from jumia on tree structure
func (p *Parser) follow(n Node) ([]*token.Token, error) {
	visited := Queue[Node]{}
	result := []*token.Token{}
	if p.isStart(n) {
		return []*token.Token{&token.Token{token.EOF, ""}}, nil
	}
	visited.enqueue(p.startNode)
	for !visited.empty() {
		curr, err := visited.dequeue()
		if err != nil {
			return nil, err
		}
		prds, err := curr.productions()
		if err != nil {
			return nil, err
		}
		for _, prd := range prds {
			// A -> bBC
			for i, ele := range prd {
				if !ele.isTerminal() {
					if ele == n {
						if i == len(prd)-1 {
							result, err = p.alphaB(result, curr)
						} else if i < len(prd)-1 {
							result, err = p.alphaBbeta(result, curr, prd, i)
						}
					}
					visited.enqueue(ele)
				}
			}
		}

	}
	return result, nil
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
