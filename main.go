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

	_ "github.com/lib/pq"
)

func main() {

	args := os.Args
	var file string = "main"
	if len(args) > 1 {
		file = args[1]
	} else {
		panic("Please specify a file")
	}

	var source string = utils.ReadFile(file + ".ev")

	tokens := lexer.Tokenize(source)

	// litter.Dump(tokens)
	// os.Exit(1)
	ast := parser.NewParser(tokens).GetAST()
	// litter.Dump(ast)
	// os.Exit(1)

	parentEnv := environment.NewEnvironment(nil)
	parentEnv.ModuleName = file
	native.SetupEnvironment(parentEnv)

	env := environment.NewEnvironment(parentEnv)
	env.ModuleName = file

	// start := time.Now()
	intr := runtime.Evaluator{Nodes: ast}
	intr.CallStack = make([]string, 0)
	intr.Evaluate(env)
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
		env.ModuleName = "repl"
		native.SetupEnvironment(env)

		run := runtime.Evaluator{Nodes: ast}
		run.Evaluate(env)
	}
}
