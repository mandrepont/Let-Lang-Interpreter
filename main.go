package main

import (
	"let_lang_proj_michael_andrepont/ast"
	"let_lang_proj_michael_andrepont/evaluator"
	"let_lang_proj_michael_andrepont/lexer"
	"let_lang_proj_michael_andrepont/parser"
	"let_lang_proj_michael_andrepont/token"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	fileName := ""
	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fmt.Print("Please input the file that contains the let program: ")
		reader := bufio.NewReader(os.Stdin)
		fileName, _ = reader.ReadString('\n')
		fileName = strings.ReplaceAll(fileName, "\n", "")
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	text := string(content)
	fmt.Println("Let Program:")
	fmt.Println(text)
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
