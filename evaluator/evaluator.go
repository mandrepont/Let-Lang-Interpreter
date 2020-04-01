package evaluator

import (
	"LetInterpreter/ast"
	"errors"
	"fmt"
)

func EvalProgram(rootNode ast.Node) (int, error) {
	if node, ok := rootNode.(ast.Expression); ok {
		return evalExpression(node, []ast.Binding{})
	} else {
		return -1, errors.New(fmt.Sprintf("Could not evaluate %T, No eval function exist for that node.", rootNode))
	}
}
func evalExpression(expressionRoot ast.Expression, e []ast.Binding) (int, error) {
	return expressionRoot.Eval(e)
}
