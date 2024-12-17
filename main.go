package main

import (
	"bufio"
	environment "evie/env"
	"evie/lexer"
	"evie/native"
	"evie/parser"
	"evie/runtime"
	"evie/utils"
	"fmt"
	"os"
)

func main() {

	args := os.Args
	var file string = "main.ev"
	if len(args) > 1 {
		file = args[1]
	} else {
		panic("Please specify a file")
	}

	var source string = utils.ReadFile(file)

	tokens := lexer.Tokenize(source)

	// litter.Dump(tokens)
	// os.Exit(1)
	ast := parser.NewParser(tokens).GetAST()
	// litter.Dump(ast[1])
	// os.Exit(1)

	parentEnv := environment.NewEnvironment(nil)
	native.SetupEnvironment(parentEnv)

	env := environment.NewEnvironment(parentEnv)
	env.ModuleName = file

	// start := time.Now()
	runtime.Evaluator{Nodes: ast}.Evaluate(env)
	// evalTime := time.Since(start).Microseconds()

	// fmt.Println("Eval time: ", evalTime, "ms")
}

func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Print(">> ")

		if !scanner.Scan() { // Lee una línea completa
			break
		}
		line := scanner.Text()

		tokens := lexer.Tokenize(line)
		ast := parser.NewParser(tokens).GetAST()

		env := environment.NewEnvironment(nil)
		native.SetupEnvironment(env)
		runtime.Evaluator{Nodes: ast}.Evaluate(env)
	}
}
