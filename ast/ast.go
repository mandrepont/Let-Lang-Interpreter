package ast

import "LetInterpreter/token"

type Node interface {
	TokenLiteral() string
}

type Expression interface {
	Node
	expressionNode()
}

type LetExpression struct {
	Token token.Token
	Name *Identifier
	Value Expression
	In Expression
}
func (le *LetExpression) expressionNode() {}
func (le *LetExpression) TokenLiteral() string { return le.Token.Literal }

type Identifier struct {
	Token token.Token //Ident Token
	Value string
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type IntLiteral struct {
	Token token.Token //INT Token
	Value int
}
func (ci *IntLiteral) expressionNode()      {}
func (ci *IntLiteral) TokenLiteral() string { return ci.Token.Literal }

type MinusExpression struct {
	Token token.Token //MINUS Token
	Arg1 Expression
	Arg2 Expression
}
func (me *MinusExpression) expressionNode() {}
func (me *MinusExpression) TokenLiteral() string { return me.Token.Literal }

type IsZeroExpression struct {
	Token token.Token //IS_ZERO Token
	Arg1 Expression
}
func (ize *IsZeroExpression) expressionNode() {}
func (ize *IsZeroExpression) TokenLiteral() string { return ize.Token.Literal }

type IfThenElseExpression struct {
	Token token.Token //IS_ZERO Token
	Predicate Expression
	TrueBranch Expression
	FalseBranch Expression
}
func (itee *IfThenElseExpression) expressionNode() {}
func (itee *IfThenElseExpression) TokenLiteral() string { return itee.Token.Literal }

