package lexer

import "github.com/medkhabt/todoprs/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Value: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	// TODO for utf-8
	var tok token.Token
	if l.ch == ':' {
		tok = newToken(token.COLON, 0)
	} else if l.ch == '/' {
		tok = newToken(token.SLASH, 0)
	} else if l.ch == '.' {
		tok = newToken(token.DOT, 0)
	} else if l.ch == ' ' {
		tok = newToken(token.SPACE, 0)
	} else if l.ch >= '0' && l.ch <= '9' {
		tok = newToken(token.DIGIT, l.ch)
	} else if l.ch == '_' || l.ch == '-' || (l.ch >= 'a' && l.ch <= 'z') || (l.ch >= 'A' && l.ch <= 'Z') {
		if l.ch == 'T' && l.peek(3) == "TODO" {
			tok = newToken(token.TODO, 0)
			l.readChar()
			l.readChar()
			l.readChar()
		} else {
			tok = newToken(token.CHAR, l.ch)
		}
	} else {
		tok = newToken(token.ILLEGAL, 0)
	}
	l.readChar()
	return tok
}
func (l *Lexer) peek(offset int) string {
	return string(l.input[l.position : l.position+offset+1])
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
