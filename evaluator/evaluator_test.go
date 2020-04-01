package evaluator

import (
	"LetInterpreter/ast"
	"fmt"
	"math"
	"strings"
	"testing"
)

func makeInt(val int) *ast.IntLiteral      { return &ast.IntLiteral{Value: val} }
func makeIdent(val string) *ast.Identifier { return &ast.Identifier{Value: val} }

func checkEvalResult(t *testing.T, expression ast.Expression, env ast.BindingList, expected int) {
	result, err := evalExpression(expression, env)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result != expected {
		t.Errorf("Expected result to be %d but was %d", expected, result)
	}
}
func checkErrorResult(t *testing.T, err error, expectedSubStr string) {
	if err == nil {
		t.Fatal("Expected eval error to exist, but it was nil.")
	} else if !strings.Contains(err.Error(), expectedSubStr) {
		t.Errorf("Expected error message to contain: [%s] but was: [%s]", expectedSubStr, err.Error())
	}
}

func TestLetEval(t *testing.T) {
	expected := 33
	expression := ast.LetExpression{
		Name:  makeIdent("y"),
		Value: makeInt(expected),
		In:    makeIdent("y"),
	}
	checkEvalResult(t, &expression, ast.BindingList{}, expected)
}

func TestIsZeroEvalTrue(t *testing.T) {
	expression := ast.IsZeroExpression{Arg1: makeInt(0)}
	checkEvalResult(t, &expression, ast.BindingList{}, 1)
}

func TestIsZeroEvalFalse(t *testing.T) {
	expression := ast.IsZeroExpression{Arg1: makeInt(10)}
	checkEvalResult(t, &expression, ast.BindingList{}, 0)
}

func TestMinusEval(t *testing.T) {
	x := []struct {
		arg1   int
		arg2   int
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
			expression := ast.MinusExpression{Arg1: makeInt(tc.arg1), Arg2: makeInt(tc.arg2)}
			checkEvalResult(t, &expression, ast.BindingList{}, tc.result)
		})
	}
}

func TestIdentBasic(t *testing.T) {
	e := ast.BindingList{
		{VarName: "x", Value: 33},
		{VarName: "test", Value: 22},
	}
	checkEvalResult(t, makeIdent("test"), e, 22)
}

func TestIdentShadowed(t *testing.T) {
	e := ast.BindingList{
		{VarName: "test", Value: 33},
		{VarName: "test", Value: 22},
	}
	checkEvalResult(t, makeIdent("test"), e, 33)
}

func TestIntLit(t *testing.T) {
	checkEvalResult(t, makeInt(0), ast.BindingList{}, 0)
	checkEvalResult(t, makeInt(math.MaxInt32), ast.BindingList{}, math.MaxInt32)
	checkEvalResult(t, makeInt(math.MinInt32), ast.BindingList{}, math.MinInt32)
}

func TestIfThenElseEvalTrue(t *testing.T) {
	expression := ast.IfThenElseExpression{
		Value:       makeInt(1),
		TrueBranch:  makeInt(22),
		FalseBranch: makeInt(33),
	}
	checkEvalResult(t, &expression, ast.BindingList{}, 22)
}

func TestIfThenElseEvalFalse(t *testing.T) {
	//Everything other than one should be false.
	expression := ast.IfThenElseExpression{
		Value:       makeInt(0),
		TrueBranch:  makeInt(22),
		FalseBranch: makeInt(33),
	}
	checkEvalResult(t, &expression, ast.BindingList{}, 33)
	expression.Value = makeInt(2)
	checkEvalResult(t, &expression, ast.BindingList{}, 33)
	expression.Value = makeInt(-1)
}

func TestIdentNotFoundEmptyEnv(t *testing.T) {
	_, err := evalExpression(makeIdent("test"), ast.BindingList{})
	checkErrorResult(t, err, "Could not find variable name: test in env of")
}

func TestIdentNotFoundNotEmptyEnv(t *testing.T) {
	e := ast.BindingList{
		{VarName: "x", Value: 33},
		{VarName: "test", Value: 22},
	}
	_, err := evalExpression(makeIdent("y"), e)
	checkErrorResult(t, err, "Could not find variable name: y in env of")
}

func TestLetInvalidValue(t *testing.T) {
	expression := ast.LetExpression{
		Name:  makeIdent("y"),
		Value: makeIdent("x"),
		In:    makeInt(7),
	}
	_, err := evalExpression(&expression, ast.BindingList{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

func TestIsZeroInvalidArg(t *testing.T) {
	expression := ast.IsZeroExpression{
		Arg1: makeIdent("x"),
	}
	_, err := evalExpression(&expression, ast.BindingList{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

func TestMinusInvalidArg1(t *testing.T) {
	expression := ast.MinusExpression{
		Arg1: makeIdent("x"),
	}
	_, err := evalExpression(&expression, ast.BindingList{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

func TestMinusInvalidArg2(t *testing.T) {
	expression := ast.MinusExpression{
		Arg1: makeIdent("x"),
		Arg2: makeIdent("y"),
	}
	_, err := evalExpression(&expression, ast.BindingList{{VarName: "x", Value: 8}})
	checkErrorResult(t, err, "Could not find variable name: y in env of")
}

func TestIfThenElseInvalidPredicate(t *testing.T) {
	expression := ast.IfThenElseExpression{
		Value: makeIdent("x"),
	}
	_, err := evalExpression(&expression, ast.BindingList{})
	checkErrorResult(t, err, "Could not find variable name: x in env of")
}

func TestAssignmentExample(t *testing.T) {
	root := ast.LetExpression{
		Name:  makeIdent("x"),
		Value: makeInt(7),
		In: &ast.LetExpression{
			Name:  makeIdent("y"),
			Value: makeInt(2),
			In: &ast.LetExpression{
				Name: makeIdent("y"),
				Value: &ast.LetExpression{
					Name: makeIdent("x"),
					Value: &ast.MinusExpression{
						Arg1: makeIdent("x"),
						Arg2: makeInt(1),
					},
					In: &ast.MinusExpression{
						Arg1: makeIdent("x"),
						Arg2: makeIdent("y"),
					},
				},
				In: &ast.MinusExpression{
					Arg1: &ast.MinusExpression{
						Arg1: makeIdent("x"),
						Arg2: makeInt(8),
					},
					Arg2: makeIdent("y"),
				},
			},
		},
	}
	result, err := EvalProgram(&root)
	expected := -5
	if err != nil {
		t.Fatal(err.Error())
	}
	if result != expected {
		t.Fatalf("Expected result to be %d but was %d", expected, result)
	}
}

func TestAssignmentExampleTwo(t *testing.T) {
	root := ast.LetExpression{
		Name:  makeIdent("x"),
		Value: makeInt(11),
		In: &ast.LetExpression{
			Name:  makeIdent("y"),
			Value: makeInt(20),
			In: &ast.IfThenElseExpression{
				Value: &ast.IsZeroExpression{
					Arg1: &ast.MinusExpression{
						Arg1: makeIdent("x"),
						Arg2: makeInt(11),
					},
				},
				TrueBranch: &ast.MinusExpression{
					Arg1: makeIdent("y"),
					Arg2: makeInt(2),
				},
				FalseBranch: &ast.MinusExpression{
					Arg1: makeIdent("y"),
					Arg2: makeInt(4),
				},
			},
		},
	}
	checkEvalResult(t, &root, ast.BindingList{}, 18)
	root.Value = makeInt(10)
	checkEvalResult(t, &root, ast.BindingList{}, 16)
}
