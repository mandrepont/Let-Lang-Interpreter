package lexer

import "let_lang_proj_michael_andrepont/token"

type Lexer struct {
	input        string
	position     int  // The current position in the input
	readPosition int  // The current reading position (one after position)
	ch           byte // The current char we are reading.
}

func New(input string) *Lexer {
	l := Lexer{input: input}
	l.readChar()
	return &l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 //EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isDigit(ch byte) bool  { return ch >= '0' && ch <= '9' }
func isLetter(ch byte) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var returnToken token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		returnToken = token.MakeToken(token.ASSIGN, l.ch)
	case ',':
		returnToken = token.MakeToken(token.COMMA, l.ch)
	case '(':
		returnToken = token.MakeToken(token.LPAREN, l.ch)
	case ')':
		returnToken = token.MakeToken(token.RPAREN, l.ch)
	case 0:
		returnToken.Type = token.EOF
		returnToken.Literal = ""
	default:
		if isLetter(l.ch) {
			returnToken.Literal = l.readIdent()
			returnToken.Type = token.KeywordLookup(returnToken.Literal)
		} else if isDigit(l.ch) {
			returnToken.Type = token.INT
			returnToken.Literal = l.readDigit()
		} else {
			returnToken = token.MakeToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return returnToken
}

func (l *Lexer) readIdent() string {
	startPos := l.position
	for isDigit(l.peekChar()) || isLetter(l.peekChar()) {
		l.readChar()
	}
	return l.input[startPos : l.position+1]
}

func (l *Lexer) readDigit() string {
	startPos := l.position
	for isDigit(l.peekChar()) {
		l.readChar()
	}
	return l.input[startPos : l.position+1]
}
