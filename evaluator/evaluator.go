package evaluator

import (
	"LetInterpreter/ast"
	"errors"
	"fmt"
)

type Binding struct {
	varName string
	value   int
}

func findIdentifierInEnv(varName string, env []Binding) (int, error) {
	for _, b := range env {
		if b.varName == varName {
			return b.value, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("Could not find variable name: %s in env of: %#v", varName, env))
}

func EvalExpression(localRoot ast.Expression, e []Binding) (int, error) {
	switch node := localRoot.(type) {
	case *ast.LetExpression:
		varName := node.Name.Value
		value, err := EvalExpression(node.Value, e)
		if err != nil { return -1, err }
		e = append([]Binding{{varName: varName, value: value}}, e...)
		return EvalExpression(node.In, e)
	case *ast.IfThenElseExpression:
		predicateVal, err := EvalExpression(node.Predicate, e)
		if err != nil { return -1, err }
		if predicateVal == 1 {
			return EvalExpression(node.TrueBranch, e)
		}
		return EvalExpression(node.FalseBranch, e)
	case *ast.IsZeroExpression:
		exprVal, err := EvalExpression(node.Arg1, e)
		if err != nil { return -1, err }
		if exprVal == 0 {
			return 1, nil
		}
		return 0, nil
	case *ast.MinusExpression:
		arg1Val, err := EvalExpression(node.Arg1, e)
		if err != nil { return -1, err }
		arg2Val, err := EvalExpression(node.Arg2, e)
		if err != nil { return -1, err }
		return arg1Val - arg2Val, nil
	case *ast.IntLiteral:
		return node.Value, nil
	case *ast.Identifier:
		return findIdentifierInEnv(node.Value, e)
	default:
		return -1, errors.New(fmt.Sprintf("Could not evaluate %T, No eval function exist for that node.", localRoot))
	}
}
