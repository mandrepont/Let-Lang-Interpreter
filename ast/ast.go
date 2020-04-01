package ast

import (
	"let_lang_proj_michael_andrepont/token"
	"errors"
	"fmt"
	"math"
	"strings"
)

type Binding struct {
	VarName string
	Value   int
}
type BindingList = []Binding

func GetEnvStr(bl *BindingList) string {
	if bl != nil {
		str := "[< "
		for _, b := range *bl {
			str += fmt.Sprintf("(%s %d) ", b.VarName, b.Value)
		}
		str += ">]"
		return str
	}
	return "[< >]"
}

func findIdentifierInEnv(varName string, env BindingList) (int, error) {
	for _, b := range env {
		if b.VarName == varName {
			return b.Value, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("Could not find variable name: %s in env of: %#v", varName, env))
}

func indentStr(indentLevel int) string {
	return strings.Repeat("\t", indentLevel)
}

//This is just extra foundation if I want to expand later.
type Node interface {
	Print(indentLevel int)
}

type Expression interface {
	Node
	Eval(env BindingList) (int, error)
	GetEnv() *BindingList
	SetEnv(*BindingList)
}

type BaseExpression struct {
	Token token.Token //IS_ZERO Token
	env   *BindingList
}

func (be *BaseExpression) GetEnv() *BindingList    { return be.env }
func (be *BaseExpression) SetEnv(env *BindingList) { be.env = env }

type LetExpression struct {
	BaseExpression
	Name  *Identifier
	Value Expression
	In    Expression
}

func (e *LetExpression) Eval(env BindingList) (int, error) {
	e.SetEnv(&env)
	e.Name.SetEnv(&env)
	varName := e.Name.Value
	value, err := e.Value.Eval(env)
	if err != nil {
		return math.MaxInt32, err
	}
	newEnv := append(BindingList{{VarName: varName, Value: value}}, env...)
	return e.In.Eval(newEnv)
}

func (e *LetExpression) Print(indentLevel int) {
	fmt.Printf("%s%s %s\n", indentStr(indentLevel), "let", GetEnvStr(e.env))
	e.Name.Print(indentLevel + 1)
	e.Value.Print(indentLevel + 1)
	e.In.Print(indentLevel + 1)
}

type Identifier struct {
	BaseExpression
	Value string
}

func (e *Identifier) Eval(env BindingList) (int, error) {
	e.SetEnv(&env)
	return findIdentifierInEnv(e.Value, env)
}
func (e *Identifier) Print(indentLevel int) {
	fmt.Printf("%s%s %s\n", indentStr(indentLevel), e.Value, GetEnvStr(e.env))
}

type IntLiteral struct {
	BaseExpression
	Value int
}

func (e *IntLiteral) Eval(env BindingList) (int, error) {
	e.SetEnv(&env)
	return e.Value, nil
}
func (e *IntLiteral) Print(indentLevel int) {
	fmt.Printf("%s%d %s\n", indentStr(indentLevel), e.Value, GetEnvStr(e.env))
}

type MinusExpression struct {
	BaseExpression
	Arg1 Expression
	Arg2 Expression
}

func (e *MinusExpression) Eval(env BindingList) (int, error) {
	e.SetEnv(&env)
	arg1Val, err := e.Arg1.Eval(env)
	if err != nil {
		return -1, err
	}
	arg2Val, err := e.Arg2.Eval(env)
	if err != nil {
		return -1, err
	}
	return arg1Val - arg2Val, nil
}
func (e *MinusExpression) Print(indentLevel int) {
	fmt.Printf("%s%s %s\n", indentStr(indentLevel), "minus", GetEnvStr(e.env))
	e.Arg1.Print(indentLevel + 1)
	e.Arg2.Print(indentLevel + 1)
}

type IsZeroExpression struct {
	BaseExpression
	Arg1 Expression
}

func (e *IsZeroExpression) Eval(env BindingList) (int, error) {
	e.SetEnv(&env)
	exprVal, err := e.Arg1.Eval(env)
	if err != nil {
		return -1, err
	}
	if exprVal == 0 {
		return 1, nil
	}
	return 0, nil
}
func (e *IsZeroExpression) Print(indentLevel int) {
	fmt.Printf("%s%s %s\n", indentStr(indentLevel), "iszero", GetEnvStr(e.env))
	e.Arg1.Print(indentLevel + 1)
}

type IfThenElseExpression struct {
	BaseExpression
	Value       Expression
	TrueBranch  Expression
	FalseBranch Expression
}

func (e *IfThenElseExpression) Eval(env BindingList) (int, error) {
	e.SetEnv(&env)
	predicateVal, err := e.Value.Eval(env)
	if err != nil {
		return -1, err
	}
	if predicateVal == 1 {
		return e.TrueBranch.Eval(env)
	}
	return e.FalseBranch.Eval(env)
}
func (e *IfThenElseExpression) Print(indentLevel int) {
	fmt.Printf("%s%s %s\n", indentStr(indentLevel), "if-then-else", GetEnvStr(e.env))
	e.Value.Print(indentLevel + 1)
	e.TrueBranch.Print(indentLevel + 1)
	e.FalseBranch.Print(indentLevel + 1)
}
