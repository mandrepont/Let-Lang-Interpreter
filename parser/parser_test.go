package parser

import (
	"let_lang_proj_michael_andrepont/ast"
	"let_lang_proj_michael_andrepont/token"
	"reflect"
	"strings"
	"testing"
)

type expressionCheck func(expression ast.Expression)

func checkForParseErrors(p *Parser, t *testing.T) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d error(s)", len(errors))
	for _, errMsg := range errors {
		t.Errorf("Parse Error: %s", errMsg)
	}
	t.FailNow()
}

func checkParseErrorsExist(p *Parser, t *testing.T, errorContains []string) {
	errors := p.Errors()
	if len(errors) != len(errorContains) {
		t.Fatalf("Expected %d parse error, but got %d", len(errorContains), len(errors))
	}
	for i, errorMsg := range errors {
		if !strings.Contains(errorMsg, errorContains[i]) {
			t.Fatalf("Expected error to contain %s, but was %s", errorContains[i], errorMsg)
		}
	}
}

func testIdent(t *testing.T, expression ast.Expression, name string) {
	v, ok := expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Parse Expression expected %T, but returned %T", &ast.Identifier{}, expression)
	}
	if v.Value != name {
		t.Fatalf("Parse Expression expected Ident name to be %s, but was %s", name, v.Value)
	}
}

func testIntLit(t *testing.T, expression ast.Expression, value int) {
	v, ok := expression.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("Parse Expression expected %T, but returned %T", &ast.Identifier{}, expression)
	}
	if v.Value != value {
		t.Fatalf("Parse Expression expected Int Lit to be %d, but was %d", value, v.Value)
	}
}

func testLetExpression(t *testing.T, expression ast.Expression, identName string, valueCheck expressionCheck, inCheck expressionCheck) {
	v, ok := expression.(*ast.LetExpression)
	if !ok {
		t.Fatalf("Parse Expression expected %T, but returned %T", &ast.LetExpression{}, expression)
	}
	testIdent(t, v.Name, identName)
	valueCheck(v.Value)
	inCheck(v.In)
}

func testMinus(t *testing.T, expression ast.Expression, arg1Check expressionCheck, arg2Check expressionCheck) {
	v, ok := expression.(*ast.MinusExpression)
	if !ok {
		t.Fatalf("Parse Expression expected %T, but returned %T", &ast.MinusExpression{}, expression)
	}
	arg1Check(v.Arg1)
	arg2Check(v.Arg2)
}

func testIsZero(t *testing.T, expression ast.Expression, argCheck expressionCheck) {
	v, ok := expression.(*ast.IsZeroExpression)
	if !ok {
		t.Fatalf("Parse Expression expected %T, but returned %T", &ast.IsZeroExpression{}, expression)
	}
	argCheck(v.Arg1)
}

func testIfThenElse(t *testing.T, expression ast.Expression, predicateCheck expressionCheck, trueCheck expressionCheck, falseCheck expressionCheck) {
	v, ok := expression.(*ast.IfThenElseExpression)
	if !ok {
		t.Fatalf("Parse Expression expected %T, but returned %T", &ast.IfThenElseExpression{}, expression)
	}
	predicateCheck(v.Value)
	trueCheck(v.TrueBranch)
	falseCheck(v.FalseBranch)
}

func TestBasicLet(t *testing.T) {
	input := []token.Token{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "8"},
		{token.IN, "in"},
		{token.IDENT, "y"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression returned nil")
	}
	testLetExpression(t, expression, "x", func(expression ast.Expression) {
		testIntLit(t, expression, 8)
	}, func(expression ast.Expression) {
		testIdent(t, expression, "y")
	})
	checkForParseErrors(p, t)
}

func TestNestedLet(t *testing.T) {
	input := []token.Token{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "8"},
		{token.IN, "in"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "9"},
		{token.IN, "in"},
		{token.IDENT, "x"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression returned nil")
	}
	testLetExpression(t, expression, "x", func(expression ast.Expression) {
		testIntLit(t, expression, 8)
	}, func(expression ast.Expression) {
		testLetExpression(t, expression, "y", func(expression ast.Expression) {
			testIntLit(t, expression, 9)
		}, func(expression ast.Expression) {
			testIdent(t, expression, "x")
		})
	})
	checkForParseErrors(p, t)
}

func TestLetMissingIdent(t *testing.T) {
	input := []token.Token{
		{token.LET, "let"},
		{token.ASSIGN, "="},
		{token.INT, "8"},
		{token.IN, "in"},
		{token.IDENT, "y"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be IDENT",
	})
}

func TestLetMissingAssign(t *testing.T) {
	input := []token.Token{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.INT, "8"},
		{token.IN, "in"},
		{token.IDENT, "y"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be =",
	})
}

func TestLetMissingValueExpr(t *testing.T) {
	input := []token.Token{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IN, "in"},
		{token.IDENT, "y"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for Value",
	})
}

func TestLetMissingInExpr(t *testing.T) {
	input := []token.Token{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "y"},
		{token.IN, "in"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for In",
	})
}

func TestLetMissingIn(t *testing.T) {
	input := []token.Token{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "y"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be IN",
	})
}

func TestIntLiteral(t *testing.T) {
	input := []token.Token{
		{token.INT, "4"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	testIntLit(t, expression, 4)
}

func TestInvalidIntLiteral(t *testing.T) {
	input := []token.Token{
		{token.INT, "let"},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Int Literal, token literal: let",
	})
}

func TestIdent(t *testing.T) {
	identLit := "testing"
	input := []token.Token{
		{token.IDENT, identLit},
		{token.EOF, ""},
	}

	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	testIdent(t, expression, identLit)
}

func TestMinus(t *testing.T) {
	input := []token.Token{
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	testMinus(t, expression, func(expression ast.Expression) {
		testIdent(t, expression, "y")
	}, func(expression ast.Expression) {
		testIntLit(t, expression, 2)
	})
}

func TestInvalidMinusMissingLParen(t *testing.T) {
	input := []token.Token{
		{token.MINUS, "minus"},
		{token.IDENT, "y"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be (",
	})
}

func TestInvalidMinusMissingRParen(t *testing.T) {
	input := []token.Token{
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.COMMA, ","},
		{token.INT, "2"},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be )",
	})
}

func TestInvalidMinusMissingComma(t *testing.T) {
	input := []token.Token{
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be ,",
	})
}

func TestInvalidMinusMissingExpression(t *testing.T) {
	input := []token.Token{
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for Arg1",
	})
}

func TestInvalidMinusMissingExpression2(t *testing.T) {
	input := []token.Token{
		{token.MINUS, "minus"},
		{token.LPAREN, "("},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.RPAREN, ")"},
	}

	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for Arg2",
	})
}

func TestIsZero(t *testing.T) {
	input := []token.Token{
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression return nil")
	}
	testIsZero(t, expression, func(expression ast.Expression) {
		testIdent(t, expression, "y")
	})
}

func TestIsZeroMissingLParen(t *testing.T) {
	input := []token.Token{
		{token.IS_ZERO, "iszero"},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be (",
	})
}

func TestIsZeroMissingRParen(t *testing.T) {
	input := []token.Token{
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be )",
	})
}

func TestIsZeroMissingArg(t *testing.T) {
	input := []token.Token{
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for Arg1",
	})
}

func TestIfThenElse(t *testing.T) {
	input := []token.Token{
		{token.IF, "if"},
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.THEN, "then"},
		{token.IDENT, "y"},
		{token.ELSE, "else"},
		{token.INT, "2"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression return nil")
	}
	testIfThenElse(t, expression, func(expression ast.Expression) {
		testIsZero(t, expression, func(expression ast.Expression) {
			testIdent(t, expression, "y")
		})
	}, func(expression ast.Expression) {
		testIdent(t, expression, "y")
	}, func(expression ast.Expression) {
		testIntLit(t, expression, 2)
	})
}

func TestIfThenElseMissingPredicate(t *testing.T) {
	input := []token.Token{
		{token.IF, "if"},
		{token.THEN, "then"},
		{token.IDENT, "y"},
		{token.ELSE, "else"},
		{token.INT, "2"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for Value",
	})
}

func TestIfThenElseMissingThen(t *testing.T) {
	input := []token.Token{
		{token.IF, "if"},
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.IDENT, "y"},
		{token.ELSE, "else"},
		{token.INT, "2"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be THEN",
	})
}

func TestIfThenElseMissingThenExpr(t *testing.T) {
	input := []token.Token{
		{token.IF, "if"},
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.THEN, "then"},
		{token.ELSE, "else"},
		{token.INT, "2"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for TrueBranch",
	})
}

func TestIfThenElseMissingElse(t *testing.T) {
	input := []token.Token{
		{token.IF, "if"},
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.THEN, "then"},
		{token.IDENT, "y"},
		{token.INT, "2"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"to be ELSE",
	})
}

func TestIfThenElseMissingElseExpr(t *testing.T) {
	input := []token.Token{
		{token.IF, "if"},
		{token.IS_ZERO, "iszero"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.THEN, "then"},
		{token.IDENT, "y"},
		{token.ELSE, "else"},
	}
	p := New(input)
	expression := p.ParseExpression()

	if !reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse Expression did not return nil")
	}
	checkParseErrorsExist(p, t, []string{
		"Missing inner expression for FalseBranch",
	})
}

func TestExample1Parse(t *testing.T) {
	input := []token.Token{
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

	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse expression should not be nil")
	}
	//This is about to be messy
	testLetExpression(t, expression, "x",
		func(e ast.Expression) {
			testIntLit(t, e, 7)
		},
		func(e ast.Expression) {
			testLetExpression(t, e, "y",
				func(e ast.Expression) {
					testIntLit(t, e, 2)
				},
				func(e ast.Expression) {
					testLetExpression(t, e, "y",
						func(e ast.Expression) {
							testLetExpression(t, e, "x",
								func(e ast.Expression) {
									testMinus(t, e, func(e ast.Expression) {
										testIdent(t, e, "x")
									}, func(e ast.Expression) {
										testIntLit(t, e, 1)
									})
								},
								func(e ast.Expression) {
									testMinus(t, e, func(e ast.Expression) {
										testIdent(t, e, "x")
									}, func(e ast.Expression) {
										testIdent(t, e, "y")
									})
								})
						},
						func(e ast.Expression) {
							testMinus(t, e, func(e ast.Expression) {
								testMinus(t, e, func(e ast.Expression) {
									testIdent(t, e, "x")
								}, func(e ast.Expression) {
									testIntLit(t, e, 8)
								})
							}, func(e ast.Expression) {
								testIdent(t, e, "y")
							})
						})
				})
		})
}

func TestExample2Parse(t *testing.T) {
	input := []token.Token{
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

	p := New(input)
	expression := p.ParseExpression()

	if reflect.ValueOf(expression).IsNil() {
		t.Fatalf("Parse expression should not be nil")
	}
	//This is about to be messy
	testLetExpression(t, expression, "x",
		func(e ast.Expression) {
			testIntLit(t, e, 11)
		},
		func(e ast.Expression) {
			testLetExpression(t, e, "y",
				func(e ast.Expression) {
					testIntLit(t, e, 20)
				},
				func(e ast.Expression) {
					testIfThenElse(t, e,
						func(e ast.Expression) {
							testIsZero(t, e, func(e ast.Expression) {
								testMinus(t, e, func(e ast.Expression) {
									testIdent(t, e, "x")
								}, func(e ast.Expression) {
									testIntLit(t, e, 11)
								})
							})
						},
						func(e ast.Expression) {
							testMinus(t, e, func(e ast.Expression) {
								testIdent(t, e, "y")
							}, func(e ast.Expression) {
								testIntLit(t, e, 2)
							})
						},
						func(e ast.Expression) {
							testMinus(t, e, func(e ast.Expression) {
								testIdent(t, e, "y")
							}, func(e ast.Expression) {
								testIntLit(t, e, 4)
							})
						})
				})
		})
}
