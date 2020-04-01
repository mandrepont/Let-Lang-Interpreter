package evaluator

import (
	"LetInterpreter/ast"
	"fmt"
	"math"
	"strings"
	"testing"
)
func makeInt(val int) *ast.IntLiteral { return &ast.IntLiteral{ Value:val } }
func makeIdent(val string) *ast.Identifier { return &ast.Identifier{ Value:val } }

func checkEvalResult(t *testing.T, expression ast.Expression, env []Binding, expected int) {
	result, err := EvalExpression(expression, env)
	if err != nil { t.Error(err.Error()) }
	if result != expected { t.Errorf("Expected result to be %d but was %d", expected, result) }
}
func checkErrorResult(t *testing.T, err error, expectedSubStr string) {
	if err == nil {
		t.Error("Expected not have a eval error.")
	} else if !strings.Contains(err.Error(), expectedSubStr) {
		t.Errorf("Expected error message to contain: [%s] but was: [%s]", expectedSubStr, err.Error())
	}
}

func TestLetEval(t *testing.T) {
	expected := 33
	expression := ast.LetExpression{
		Name: makeIdent("y"),
		Value: makeInt(expected),
		In: makeIdent("y"),
	}
	checkEvalResult(t, &expression, []Binding{}, expected)
}

func TestIsZeroEvalTrue(t *testing.T) {
	expression := ast.IsZeroExpression{ Arg1:  makeInt(0), }
	checkEvalResult(t, &expression, []Binding{}, 1)
}

func TestIsZeroEvalFalse(t *testing.T) {
	expression := ast.IsZeroExpression{ Arg1:  makeInt(10), }
	checkEvalResult(t, &expression, []Binding{}, 0)
}

func TestMinusEval(t *testing.T) {
	x := []struct {
		arg1 int
		arg2 int
		result int
	}{
		{1, 3, -2},
		{-2, -5, 3},
		{10, 10, 0},
		{math.MaxInt32, math.MaxInt32, 0},
		{math.MinInt32, math.MinInt32, 0},
	}
	for _, tc := range x {
		t.Run(fmt.Sprintf("%d-%d=%d", tc.arg1, tc.arg2, tc.result), func(t *testing.T) {
			expression := ast.MinusExpression{ Arg1:  makeInt(tc.arg1), Arg2:makeInt(tc.arg2) }
			checkEvalResult(t, &expression, []Binding{}, tc.result)
		})
	}
}

func TestIdentBasic(t *testing.T) {
	e := []Binding{
		{varName: "x", value: 33},
		{varName: "test", value: 22},
	}
	checkEvalResult(t, makeIdent("test"), e, 22)
}

func TestIdentShadowed(t *testing.T) {
	e := []Binding{
		{varName: "test", value: 33},
		{varName: "test", value: 22},
	}
	checkEvalResult(t, makeIdent("test"), e, 33)
}

func TestIntLit(t *testing.T) {
	checkEvalResult(t, makeInt(0), []Binding{}, 0)
	checkEvalResult(t, makeInt(math.MaxInt32), []Binding{}, math.MaxInt32)
	checkEvalResult(t, makeInt(math.MinInt32), []Binding{}, math.MinInt32)
}

func TestIfThenElseEvalTrue(t *testing.T) {
	expression := ast.IfThenElseExpression{
		Predicate:   makeInt(1),
		TrueBranch:  makeInt(22),
		FalseBranch: makeInt(33),
	}
	checkEvalResult(t, &expression, []Binding{}, 22)
}

func TestIfThenElseEvalFalse(t *testing.T) {
	//Everything other than one should be false.
	expression := ast.IfThenElseExpression{
		Predicate:   makeInt(0),
		TrueBranch:  makeInt(22),
		FalseBranch: makeInt(33),
	}
	checkEvalResult(t, &expression, []Binding{}, 33)
	expression.Predicate = makeInt(2)
	checkEvalResult(t, &expression, []Binding{}, 33)
	expression.Predicate = makeInt(-1)
}

// Error testing
func TestExpressionNotSupported(t *testing.T) {
	_, err := EvalExpression(&ast.TestNode{Value: ""}, []Binding{})
	checkErrorResult(t, err, "Could not evaluate *ast.TestNode")
}

func TestIdentNotFoundEmptyEnv(t *testing.T) {
	_, err := EvalExpression(makeIdent("test"), []Binding{})
	checkErrorResult(t, err, "Could not find variable name: test in env of")
}

func TestIdentNotFoundNotEmptyEnv(t *testing.T) {
	e := []Binding{
		{varName: "x", value: 33},
		{varName: "test", value: 22},
	}
	_, err := EvalExpression(makeIdent("y"), e)
	checkErrorResult(t, err, "Could not find variable name: y in env of")
}

func TestLetInvalidValue(t *testing.T) {
	expression := ast.LetExpression{
		Name:  makeIdent("y"),
		Value: makeIdent("x"),
		In:    makeInt(7),
	}
	_, err := EvalExpression(&expression, []Binding{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

func TestIsZeroInvalidArg(t *testing.T) {
	expression := ast.IsZeroExpression{
		Arg1:  makeIdent("x"),
	}
	_, err := EvalExpression(&expression, []Binding{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

func TestMinusInvalidArg1(t *testing.T) {
	expression := ast.MinusExpression{
		Arg1:  makeIdent("x"),
	}
	_, err := EvalExpression(&expression, []Binding{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

func TestMinusInvalidArg2(t *testing.T) {
	expression := ast.MinusExpression{
		Arg1:  makeIdent("x"),
		Arg2:  makeIdent("y"),
	}
	_, err := EvalExpression(&expression, []Binding{{varName: "x", value: 8}})
	checkErrorResult(t, err, "Could not find variable name: y in env of")
}

func TestIfThenElseInvalidPredicate(t *testing.T) {
	expression := ast.IfThenElseExpression{
		Predicate: makeIdent("x"),
	}
	_, err := EvalExpression(&expression, []Binding{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

