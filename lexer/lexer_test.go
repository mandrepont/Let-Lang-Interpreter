package lexer

import (
	"LetInterpreter/token"
	"testing"
)

type ExpectedTokens []token.Token

func checkTokens(t *testing.T, input string, expectedToken ExpectedTokens) {
	lexer := New(input)
	for i, et := range expectedToken {
		nextToken := lexer.NextToken()

		if nextToken != et {
			t.Fatalf("nextToken[%d] - nextToken is not expected. expected=%+v, got=%+v",
				i, et, nextToken)
		}
	}
}

func TestSingleTokenLex(t *testing.T) {
	input := `=(),`
	expectedTokens := ExpectedTokens{
		{token.ASSIGN, "="},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

func TestSingleTokenWithIllegalLex(t *testing.T) {
	input := `=(),[]{}`
	expectedTokens := ExpectedTokens{
		{token.ASSIGN, "="},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.ILLEGAL, "["},
		{token.ILLEGAL, "]"},
		{token.ILLEGAL, "{"},
		{token.ILLEGAL, "}"},
		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

func TestDigitsAndIdentLex(t *testing.T) {
	input := `myAwesome83StRIng 321 32n2x 3x3 3 y`
	expectedTokens := ExpectedTokens{
		{token.IDENT, "myAwesome83StRIng"},
		{token.INT, "321"},
		{token.INT, "32"},
		{token.IDENT, "n2x"},
		{token.INT, "3"},
		{token.IDENT, "x3"},
		{token.INT, "3"},
		{token.IDENT, "y"},
		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

func TestKeywordsLex(t *testing.T) {
	input := `let iszero mincus minus if then else in`
	expectedTokens := ExpectedTokens{
		{token.LET, "let"},
		{token.IS_ZERO, "iszero"},
		{token.IDENT, "mincus"},
		{token.MINUS, "minus"},
		{token.IF, "if"},
		{token.THEN, "then"},
		{token.ELSE, "else"},
		{token.IN, "in"},
		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

func TestAssignmentExample(t *testing.T) {
	input := `
		let x = 7
		in let y = 2
			in let y = let x = minus(x, 1)
				in minus(x, y)
			in minus(minus(x, 8), y)
	`
	expectedTokens := ExpectedTokens{
		//Break = new line
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "7"},

		{token.IN, "in"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "2"},

		{token.IN, "in"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.INT, "1"},
		{token.RPAREN, ")"},

		{token.IN, "in"},
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},

		{token.IN, "in"},
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.INT, "8"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},

		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

func TestAssignment2Example(t *testing.T) {
	input := `
		let x = 11
		in let y = 20
			in if iszero(minus(x, 11)) 
				then minus(y, 2)
				else minus(y, 4)
	`
	expectedTokens := ExpectedTokens{
		//Break = new line
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "11"},

		{token.IN, "in"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "20"},

		{token.IN, "in"},
		{token.IF, "if"},
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.INT, "11"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},

		{token.THEN, "then"},
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},

		{token.ELSE, "else"},
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.COMMA, ","},
		{token.INT, "4"},
		{token.RPAREN, ")"},

		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

func TestOnlyEOF(t *testing.T) {
	input := ``
	expectedTokens := ExpectedTokens{
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

func TestEOFWithToken(t *testing.T) {
	input := `let x = 8`
	expectedTokens := ExpectedTokens{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "8"},

		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
		{token.EOF, ""},
	}
	checkTokens(t, input, expectedTokens)
}

