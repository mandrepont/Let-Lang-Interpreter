package ast

import (
	"LetInterpreter/token"
	"errors"
	"fmt"
	"math"
)

type Binding struct {
	VarName string
	Value   int
}

func findIdentifierInEnv(varName string, env []Binding) (int, error) {
	for _, b := range env {
		if b.VarName == varName {
			return b.Value, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("Could not find variable name: %s in env of: %#v", varName, env))
}

//This is just extra foundation if I want to expand later.
type Node interface {
}

type Expression interface {
	Node
	Eval(env []Binding) (int, error)
}

type LetExpression struct {
	Token token.Token
	Name  *Identifier
	Value Expression
	In    Expression
}

func (le *LetExpression) Eval(env []Binding) (int, error) {
	varName := le.Name.Value
	value, err := le.Value.Eval(env)
	if err != nil {
		return math.MaxInt32, err
	}
	env = append([]Binding{{VarName: varName, Value: value}}, env...)
	return le.In.Eval(env)
}

type Identifier struct {
	Token token.Token //Ident Token
	Value string
}

func (i *Identifier) Eval(env []Binding) (int, error) { return findIdentifierInEnv(i.Value, env) }

type IntLiteral struct {
	Token token.Token //INT Token
	Value int
}

func (i *IntLiteral) Eval(_ []Binding) (int, error) { return i.Value, nil }

type MinusExpression struct {
	Token token.Token //MINUS Token
	Arg1  Expression
	Arg2  Expression
}

func (mExpr *MinusExpression) Eval(env []Binding) (int, error) {
	arg1Val, err := mExpr.Arg1.Eval(env)
	if err != nil {
		return -1, err
	}
	arg2Val, err := mExpr.Arg2.Eval(env)
	if err != nil {
		return -1, err
	}
	return arg1Val - arg2Val, nil
}

type IsZeroExpression struct {
	Token token.Token //IS_ZERO Token
	Arg1  Expression
}

func (iz *IsZeroExpression) Eval(env []Binding) (int, error) {
	exprVal, err := iz.Arg1.Eval(env)
	if err != nil {
		return -1, err
	}
	if exprVal == 0 {
		return 1, nil
	}
	return 0, nil
}

type IfThenElseExpression struct {
	Token       token.Token //IS_ZERO Token
	Value       Expression
	TrueBranch  Expression
	FalseBranch Expression
}

func (itee *IfThenElseExpression) Eval(env []Binding) (int, error) {
	predicateVal, err := itee.Value.Eval(env)
	if err != nil {
		return -1, err
	}
	if predicateVal == 1 {
		return itee.TrueBranch.Eval(env)
	}
	return itee.FalseBranch.Eval(env)
}
