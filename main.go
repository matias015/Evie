package main

import (
	environment "evie/env"
	"evie/evruntime"
	"evie/lexer"
	"evie/native"
	"evie/parser"
	"evie/utils"
	"os"
)

func main() {

	file := GetFileName()

	var source string = utils.ReadFile(utils.AddExtension(file))

	tokens := lexer.Tokenize(source)

	// litter.Dump(tokens)
	// os.Exit(1)
	ast := parser.NewParser(tokens).GetAST()
	// litter.Dump(ast)
	// os.Exit(1)

	var env *environment.Environment = environment.NewEnvironment()

	native.SetupEnvironment(env)

	env.ModuleName = file
	env.ImportChain[file] = true

	// start := time.Now()
	intr := evruntime.Evaluator{Nodes: ast}
	intr.Evaluate(env)
	// evalTime := time.Since(start).Microseconds()

	// fmt.Println("Eval time: ", evalTime, "ms")
}

func GetFileName() string {
	args := os.Args

	var file string = "main"

	if len(args) > 1 {
		file = args[1]
	} else {
		panic("Please specify a file")
	}

	return file
}
