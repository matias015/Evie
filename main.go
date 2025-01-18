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
	"path"
	"path/filepath"
	"time"
)

func Init() {

	// timer := profil.ObtenerInstancia()

	file := GetFileName()

	ast := ParseContent(file)

	// GetRootPath(file)

	var env *environment.Environment = SetupInitialEnv(file)

	start := time.Now()

	intr := evruntime.Evaluator{Nodes: ast}
	intr.RootPath = GetRootPath(file)
	intr.Evaluate(env)

	fmt.Println("\nEval time: ", time.Since(start).Microseconds()/1000, "ms")
	// timer.Display()
}

func GetRootPath(mainFileName string) string {

	cd, err := os.Getwd()

	if err != nil {
		fmt.Println("Error trying to get the actual working directory")
	}

	root := path.Dir(cd + string(filepath.Separator) + common.AddExtension(mainFileName))

	// err = os.Chdir(root)

	// if err != nil {
	// 	fmt.Println("Error setting up the root directory")
	// 	os.Exit(1)
	// }

	// cd, _ = os.Getwd()

	return root
}

func ParseContent(file string) []parser.Stmt {
	var source string = common.ReadFile(common.AddExtension(file))

	tokens := lexer.Tokenize(source)

	ast := parser.NewParser(tokens).GetAST()

	return ast
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
