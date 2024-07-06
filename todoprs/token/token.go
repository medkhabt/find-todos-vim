package token

const (
	EOF     = "EOF"
	DOT     = "."
	COLON   = ":"
	DIGIT   = "DIGIT"
	SLASH   = "/"
	CHAR    = "char"
	TODO    = "TODO"
	SPACE   = "SPACE"
	ILLEGAL = "ILLEGAL"
	EPSILON = "EPSILON"
	//COMMENTSIGN = "COMMENTSIGN"
)

type TokenType string
type Token struct {
	Type  TokenType
	Value string
}

func ExistEpsTk(tks []*Token) bool {
	for _, tok := range tks {
		if tok.Type == EPSILON {
			return true
		}
	}
	return false
}
