package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func MakeToken(tokenType TokenType, char byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(char),
	}
}

func KeywordLookup(literal string) TokenType {
	keywordsMap := map[string]TokenType{
		"if":     IF,
		"else":   ELSE,
		"then":   THEN,
		"let":    LET,
		"in":     IN,
		"minus":  MINUS,
		"iszero": IS_ZERO,
	}
	if tokType, ok := keywordsMap[literal]; ok {
		return tokType
	}
	return IDENT
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//Ident and lit
	IDENT = "IDENT"
	INT   = "INT"

	//Keywords
	LET     = "LET"
	IN      = "IN"
	IF      = "IF"
	THEN    = "THEN"
	ELSE    = "ELSE"
	IS_ZERO = "IS_ZERO"
	MINUS   = "MINUS"

	ASSIGN = "="
	COMMA  = ","
	LPAREN = "("
	RPAREN = ")"
)
