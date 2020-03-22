package parser

import (
	"LetInterpreter/ast"
	"LetInterpreter/token"
	"fmt"
	"strconv"
)

type Parser struct {
	tokenQueue []token.Token
	currentToken token.Token
	peekToken token.Token
	position int
	errors []string
}

func New(tokenQueue []token.Token) *Parser {
	p := &Parser{tokenQueue: tokenQueue, position: -1, peekToken: tokenQueue[0]}
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.position++
	p.currentToken = p.peekToken
	if p.position + 1 >= len(p.tokenQueue) {
		p.peekToken = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	} else {
		p.peekToken = p.tokenQueue[p.position+1]
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool { return p.currentToken.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }
func (p *Parser) peekError(t token.TokenType) {
	errorMsg := fmt.Sprintf("Excpected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, errorMsg)
}
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) ParseExpression() ast.Expression {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetExpression()
	case token.IDENT:
		return p.parseIdentifier()
	case token.INT:
		return p.parseIntLiteral()
	case token.MINUS:
		return p.parseMinusExpression()
	case token.IS_ZERO:
		return p.parseIsZeroExpression()
	case token.IF:
		return p.parseIfThenElseExpression()
	}
	return nil
}

func (p *Parser) parseLetExpression() *ast.LetExpression {
	expr := &ast.LetExpression{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) { return nil }

	expr.Name = p.parseIdentifier()
	if expr.Name == nil { return nil }

	if !p.expectPeek(token.ASSIGN) { return nil }

	p.nextToken()
	expr.Value = p.ParseExpression()
	if expr.Value == nil {
		p.errors = append(p.errors, "Missing inner expression for Value")
		return nil
	}

	if !p.expectPeek(token.IN) { return nil }

	p.nextToken()
	expr.In = p.ParseExpression()
	if expr.In == nil {
		p.errors = append(p.errors, "Missing inner expression for In")
		return nil
	}

	return expr
}

func (p *Parser) parseMinusExpression() *ast.MinusExpression {
	expr := &ast.MinusExpression{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) { return nil }

	p.nextToken()
	expr.Arg1 = p.ParseExpression()
	if expr.Arg1 == nil {
		p.errors = append(p.errors, "Missing inner expression for Arg1")
		return nil
	}

	if !p.expectPeek(token.COMMA) { return nil }

	p.nextToken()
	expr.Arg2 = p.ParseExpression()
	if expr.Arg2 == nil {
		p.errors = append(p.errors, "Missing inner expression for Arg2")
		return nil
	}

	if !p.expectPeek(token.RPAREN) { return nil }
	return expr
}

func (p *Parser) parseIsZeroExpression() *ast.IsZeroExpression {
	expr := &ast.IsZeroExpression{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) { return nil }

	p.nextToken()
	expr.Arg1 = p.ParseExpression()
	if expr.Arg1 == nil {
		p.errors = append(p.errors, "Missing inner expression for Arg1")
		return nil
	}

	if !p.expectPeek(token.RPAREN) { return nil }
	return expr
}

func (p *Parser) parseIfThenElseExpression() *ast.IfThenElseExpression {
	expr := &ast.IfThenElseExpression{Token: p.currentToken}

	p.nextToken()
	expr.Predicate = p.ParseExpression()
	if expr.Predicate == nil {
		p.errors = append(p.errors, "Missing inner expression for Predicate")
		return nil
	}

	if !p.expectPeek(token.THEN) { return nil }

	p.nextToken()
	expr.TrueBranch = p.ParseExpression()
	if expr.TrueBranch == nil {
		p.errors = append(p.errors, "Missing inner expression for TrueBranch")
		return nil
	}

	if !p.expectPeek(token.ELSE) { return nil }

	p.nextToken()
	expr.FalseBranch = p.ParseExpression()
	if expr.FalseBranch == nil {
		p.errors = append(p.errors, "Missing inner expression for FalseBranch")
		return nil
	}

	return expr
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	ident := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	return ident
}

func (p *Parser) parseIntLiteral() *ast.IntLiteral {
	value, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("Error parsing Int Literal, token literal: %s, Atio error: %s",
			p.currentToken.Literal, err.Error())
		p.errors = append(p.errors, msg)
		return nil
	}
	intLit := &ast.IntLiteral{
		Token: p.currentToken,
		Value: value,
	}
	return intLit
}
