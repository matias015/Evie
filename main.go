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

	args := os.Args
	var file string = "main"
	if len(args) > 1 {
		file = args[1]
	} else {
		panic("Please specify a file")
	}

	var source string = utils.ReadFile(file + ".ev")
	// for _, char := range source {
	// 	fmt.Println(char)
	// }
	// os.Exit(1)

	tokens := lexer.Tokenize(source)

	// litter.Dump(tokens)
	// os.Exit(1)
	ast := parser.NewParser(tokens).GetAST()
	// litter.Dump(ast[1])
	// os.Exit(1)

	parentEnv := environment.NewEnvironment(nil)
	parentEnv.ModuleName = file
	native.SetupEnvironment(parentEnv)

	env := environment.NewEnvironment(parentEnv)
	env.ModuleName = file

	// start := time.Now()
	intr := evruntime.Evaluator{Nodes: ast}
	intr.CallStack = make([]string, 0)

	intr.Evaluate(env)
	// evalTime := time.Since(start).Microseconds()

	// fmt.Println("Eval time: ", evalTime, "ms")
}

func Test() {

}
