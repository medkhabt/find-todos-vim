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
	String() string
}

type Nodes []Node

func (nodes Nodes) String() string {
	s := "["
	for i, node := range nodes {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", node.String())
	}
	return s + "]"
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
	black := make(map[string]bool)
	grey := Queue[Node]{}
	result := []*token.Token{}
	if p.isStart(n) {
		return []*token.Token{&token.Token{token.EOF, ""}}, nil
	}
	grey.enqueue(p.startNode)
	for !grey.empty() {
		t, err := grey.dequeue()
		curr := t.(*NonTerminalNode)
		black[curr.name] = true
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
					ele := ele.(*NonTerminalNode)
					if !black[ele.name] {
						grey.enqueue(ele)
					}
				}
			}
		}

	}
	return result, nil
}

// for each A -> \alpha
//
//	terminal a in FIRST(\alph), add A -> \alph to M[A,a]
//	if eps in FIRST(\alph), then foreach b terminal in FOLLOW(A), add A -> \alpha to M[A, $]
type Transition struct {
	nodeName      string
	transitorType token.TokenType
}

// LL1, so we have one production per key.
func (p *Parser) makeParsingTable() (map[Transition][]Node, error) {
	black := make(map[string]bool)
	grey := Queue[Node]{}
	// example A-> \alpha to (A, a)
	m := make(map[Transition][]Node)

	grey.enqueue(p.startNode)
	for !grey.empty() {
		t, err := grey.dequeue()
		if err != nil {
			return nil, err
		}
		curr, ok := t.(*NonTerminalNode)
		black[curr.name] = true
		if !ok {
			return nil, fmt.Errorf("curr is not a NonTerminlaNode")
		}
		prds, err := curr.productions()
		if err != nil {
			return nil, err
		}
		for _, prd := range prds {
			// A -> bBC
			alpha := (&NonTerminalNode{}).addProduction(prd)
			for _, ele := range prd {
				if !ele.isTerminal() {
					name := ele.(*NonTerminalNode).name
					if !black[name] {
						grey.enqueue(ele)
					}
				}
			}
			firsts, err := p.first(alpha)
			if err != nil {
				return nil, err
			}
			epsilonExist := false
			for _, f := range firsts {
				if f.Type == token.EPSILON {
					epsilonExist = true
				} else {
					// 1) For each terminal a in FIRST(\alpha), add A -> \alpha to M[A,a]
					m[Transition{curr.name, f.Type}] = prd
				}
			}
			// if \epsilon exists in FIRST(\alpha)
			if epsilonExist {
				follows, err := p.follow(curr)
				if err != nil {
					return nil, err
				}
				for _, fo := range follows {
					// 2) For each terminal b in FOLLOW(A), add A -> \alpha to M[A,b] (EOF included)
					m[Transition{curr.name, fo.Type}] = prd
				}
			}
		}
	}
	return m, nil
}

// TODO put the stack ds and the queue in other package.
type Stack[E any] interface {
	Push(e E) error
	Pop() (E, error)
	Peek() (E, error)
}
type BasicStack[E any] struct {
	list []E
	top  int
}

func NewBasicStack[E any](capacity int) *BasicStack[E] {
	return &BasicStack[E]{make([]E, capacity, capacity), -1}
}

func (s *BasicStack[E]) Push(e E) error {
	if s.top == cap(s.list) {
		return fmt.Errorf("Exceding the capacity of the Stack.")
	}
	s.top += 1
	s.list[s.top] = e
	return nil
}
func (s *BasicStack[E]) Pop() (E, error) {
	if s.top < 0 {
		var zero E
		return zero, fmt.Errorf("Stack is empty, can't pop.")
	}
	r := s.list[s.top]
	s.top -= 1
	return r, nil
}

func (s *BasicStack[E]) Peek() (E, error) {
	if s.top < 0 {
		var zero E
		return zero, fmt.Errorf("Stack is empty, can't pop.")
	}
	return s.list[s.top], nil
}

// TODO check if input has EOF token , should it have one ?
func (p *Parser) PredictiveParsing(inputBuffer []*token.Token, prsTbl map[Transition][]Node) error {
	counter := 0
	// random size for now
	var stack Stack[Node] = NewBasicStack[Node](100)
	var dollar Node = &TerminalNode{&token.Token{token.EOF, ""}}
	// TODO should I implement the stack with pointers as args?
	err := stack.Push(dollar)
	if err != nil {
		return err
	}
	stack.Push(p.startNode)
	if err != nil {
		return err
	}
	a := inputBuffer[counter]
	X, err := stack.Peek()
	if err != nil {
		return err
	}

	for X != dollar {
		X, err = stack.Peek()
		if err != nil {
			return err
		}
		if X.isTerminal() {
			Y, ok := X.(*TerminalNode)
			if !ok {
				fmt.Errorf("Couldn't type assert %v to TerminalNode", X)
			}
			if *(Y.token) == *a {
				stack.Pop()
				counter++
				a = inputBuffer[counter]
			} else {
				// parsing error
				return fmt.Errorf("Parsing error for terminal")
			}
		} else {
			Y, ok := X.(*NonTerminalNode)
			if !ok {
				fmt.Errorf("Couldn't type assert %v to *NonTerminalNode", X)
			}
			prod, ok := prsTbl[Transition{Y.name, a.Type}]
			stack.Pop()
			if ok {
				fmt.Printf("%s -> %v \n", Y.name, prod)
				N := len(prod)
				for i := N - 1; i >= 0; i -= 1 {
					stack.Push(prod[i])
				}
			} else {
				return fmt.Errorf("Parsing error for non terminal")
				// parsing error
			}
		}
		X, err = stack.Peek()
	}
	if err != nil {
		return err
	}
	return nil
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
