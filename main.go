package main

import (
	"evie/common"
	environment "evie/env"
	"evie/evruntime"
	"evie/lexer"
	"evie/native"
	"evie/parser"
	"fmt"
	"os"
	"time"
)

func Init() {

	// timer := profil.ObtenerInstancia()

	file := GetFileName()

	ast := ParseContent(file)

	var env *environment.Environment = SetupInitialEnv(file)

	start := time.Now()

	intr := evruntime.Evaluator{Nodes: ast}
	intr.Evaluate(env)

	fmt.Println("\nEval time: ", time.Since(start).Microseconds()/1000, "ms")
	// timer.Display()
}

func ParseContent(file string) []parser.Stmt {
	var source string = common.ReadFile(common.AddExtension(file))

	tokens := lexer.Tokenize(source)

	return parser.NewParser(tokens).GetAST()
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

func SetupInitialEnv(moduleName string) *environment.Environment {
	var env *environment.Environment = environment.NewEnvironment()

	native.SetupEnvironment(env)

	env.ModuleName = moduleName
	env.ImportChain[moduleName] = true
	return env
}

func main() {

	// Parse cl arguments

	Init()
}
