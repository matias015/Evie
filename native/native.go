package native

import (
	"bufio"
	environment "evie/env"
	"evie/values"
	"fmt"
	"os"

	"time"

	"github.com/sanity-io/litter"
)

func SetupEnvironment(env *environment.Environment) {

	env.Variables["input"] = values.NativeFunctionValue{Value: ReadUserInput}

	env.Variables["print"] = values.NativeFunctionValue{Value: PrintStdOut}

	env.DeclareVar("time", values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		now := time.Now()
		milliseconds := now.UnixMilli()
		return values.NumberValue{Value: float64(milliseconds)}
	}})

	env.DeclareVar("litter", values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		litter.Dump(args[0])
		return values.BooleanValue{Value: true}
	}})

	env.DeclareVar("panic", values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		return values.ErrorValue{Value: "Panic! -> " + args[0].GetStr()}
	}})

	arguments := []values.RuntimeValue{}

	for _, arg := range os.Args {
		arguments = append(arguments, values.StringValue{Value: arg})
	}

	env.Variables["getArgs"] = values.NativeFunctionValue{Value: GetArguments}
}

func ReadUserInput(args []values.RuntimeValue) values.RuntimeValue {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() { // Lee una línea completa
		return values.StringValue{Value: ""}
	}
	line := scanner.Text()
	return values.StringValue{Value: line}
}

func PrintStdOut(args []values.RuntimeValue) values.RuntimeValue {

	for _, arg := range args {
		if arg.GetType() == "ArrayValue" {
			fmt.Print("[ ")
			for _, item := range arg.(values.ArrayValue).Value {
				PrintStdOut([]values.RuntimeValue{item})
				fmt.Print(", ")
			}
			fmt.Print("] ")
		}
		if arg.GetType() == "StringValue" {
			fmt.Print("\"" + arg.GetStr() + "\"")
		} else if arg.GetType() == "NumberValue" {
			fmt.Print(arg.GetNumber())
		} else if arg.GetType() == "BooleanValue" {
			fmt.Print(arg.GetBool())
		} else if arg.GetType() == "DictionaryValue" {
			fmt.Print("{ ")
			for key, value := range arg.(values.DictionaryValue).Value {
				PrintStdOut([]values.RuntimeValue{values.StringValue{Value: key}})
				fmt.Print(": ")
				PrintStdOut([]values.RuntimeValue{value})
				fmt.Print(", ")
			}
			fmt.Print("}")
		}
	}
	print("\n")

	return values.BooleanValue{Value: true}
}

func GetArguments(args []values.RuntimeValue) values.RuntimeValue {

	vals := values.ArrayValue{Value: []values.RuntimeValue{}}

	sysargs := os.Args

	for _, arg := range sysargs {
		vals.Value = append(vals.Value, values.StringValue{Value: arg})
	}

	return vals

}
