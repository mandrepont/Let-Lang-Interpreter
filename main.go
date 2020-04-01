package main

import (
	"LetInterpreter/ast"
	"LetInterpreter/evaluator"
	"LetInterpreter/lexer"
	"LetInterpreter/parser"
	"LetInterpreter/token"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	fileName := ""
	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		println("Usage:\n go run main.go 'file_to_eval.let'")
		return
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	text := string(content)
	tokens := getTokenList(text)
	root := getAst(tokens)
	printEvalResult(root)
}

func getTokenList(lexerInput string) []token.Token {
	lxr := lexer.New(lexerInput)
	var tokens []token.Token
	lexedToken := lxr.NextToken()
	for lexedToken.Type != token.EOF {
		tokens = append(tokens, lexedToken)
		lexedToken = lxr.NextToken()
	}
	tokens = append(tokens, lexedToken)
	fmt.Println("\nToken Queue:")
	for _, t := range tokens {
		fmt.Printf("%+v\n", t)
	}
	return tokens
}

func getAst(tokens []token.Token) ast.Node {
	prs := parser.New(tokens)
	expr := prs.ParseExpression()
	if len(prs.Errors()) > 0 {
		for _, err := range prs.Errors() {
			log.Fatalf(err)
		}
	}
	fmt.Println("\nAST without env:")
	expr.Print(0)
	return expr
}

func printEvalResult(root ast.Node) {
	res, err := evaluator.EvalProgram(root)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("\nAST with env:")
	root.Print(0)
	fmt.Println("\nExpression Result: ", res)
}
