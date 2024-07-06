package lexer

import (
	"testing"

	"github.com/medkhabt/todoprs/token"
)

type TokenTest struct {
	expectedType  token.TokenType
	expectedValue string
}

func TestNextToken(t *testing.T) {
	input := `t/ed/a.d:2://TODO t`
	tests := []*TokenTest{
		&TokenTest{token.CHAR, "t"},
		&TokenTest{token.SLASH, "\x00"},
		&TokenTest{token.CHAR, "e"},
		&TokenTest{token.CHAR, "d"},
		&TokenTest{token.SLASH, "\x00"},
		&TokenTest{token.CHAR, "a"},
		&TokenTest{token.DOT, "\x00"},
		&TokenTest{token.CHAR, "d"},
		&TokenTest{token.COLON, "\x00"},
		&TokenTest{token.DIGIT, "2"},
		&TokenTest{token.COLON, "\x00"},
		&TokenTest{token.SLASH, "\x00"},
		&TokenTest{token.SLASH, "\x00"},
		&TokenTest{token.TODO, "\x00"},
		&TokenTest{token.SPACE, "\x00"},
		&TokenTest{token.CHAR, "t"},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		t.Logf("tok is %v.\n", tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d]: token type don't match. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Value != tt.expectedValue {
			t.Fatalf("test[%d]: token value don't match. expected=%q, got=%q", i, tt.expectedValue, tok.Value)
		}
	}
}
