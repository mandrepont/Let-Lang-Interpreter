package ast

import "LetInterpreter/token"

//This is just extra foundation if I want to expand later.
type Node interface {
}

type Expression interface {
	Node
}

type LetExpression struct {
	Token token.Token
	Name *Identifier
	Value Expression
	In Expression
}
func (le *LetExpression) eval() int {

}

type Identifier struct {
	Token token.Token //Ident Token
	Value string
}

type IntLiteral struct {
	Token token.Token //INT Token
	Value int
}

type MinusExpression struct {
	Token token.Token //MINUS Token
	Arg1 Expression
	Arg2 Expression
}

type IsZeroExpression struct {
	Token token.Token //IS_ZERO Token
	Arg1 Expression
}

type IfThenElseExpression struct {
	Token token.Token //IS_ZERO Token
	Predicate Expression
	TrueBranch Expression
	FalseBranch Expression
}

//DO NOT USE: Node used for unit testing for invalid test case.
type TestNode struct {
	Value string
}

