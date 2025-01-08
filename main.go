package main

import (
	environment "evie/env"
	"evie/evruntime"
	"evie/lexer"
	"evie/native"
	"evie/parser"
	"evie/utils"
	"fmt"
	"os"
	"time"
)

func main() {

	// timer := profil.ObtenerInstancia()

	file := GetFileName()

	var source string = utils.ReadFile(utils.AddExtension(file))

	tokens := lexer.Tokenize(source)

	ast := parser.NewParser(tokens).GetAST()

	var env *environment.Environment = SetupInitialEnv(file)

	start := time.Now()

	intr := evruntime.Evaluator{Nodes: ast}
	intr.Evaluate(env)

	fmt.Println("\nEval time: ", time.Since(start).Microseconds()/1000, "ms")
	// timer.Display()

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

func SetupInitialEnv(file string) *environment.Environment {
	var env *environment.Environment = environment.NewEnvironment()

	native.SetupEnvironment(env)

	env.ModuleName = file
	env.ImportChain[file] = true
	return env
}
