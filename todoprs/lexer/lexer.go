package lexer

import "github.com/medkhabt/todoprs/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	// I don't know if this a good solution but it solves my issue with whitespaces, i can keep the ones after TODO.
	inTask bool
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
		if l.inTask {
			tok = newToken(token.SPACE, 0)
		} else {
			l.readChar()
			tok = l.NextToken()
		}
	} else if l.ch == byte(0x0D) || l.ch == byte(0x0A) {
		tok = newToken(token.NEWLINE, 0)
		s := l.peek(1)
		if s == "\n\r" || s == "\r\n" {
			l.readChar()
		}
		l.inTask = false
	} else if l.ch >= '0' && l.ch <= '9' {
		tok = newToken(token.DIGIT, l.ch)
	} else if l.ch == '_' || l.ch == '-' || (l.ch >= 'a' && l.ch <= 'z') || (l.ch >= 'A' && l.ch <= 'Z') {
		if l.ch == 'T' && l.peek(3) == "TODO" {
			tok = newToken(token.TODO, 0)
			l.readChar()
			l.readChar()
			l.readChar()
			l.inTask = true
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
