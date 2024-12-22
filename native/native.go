package native

import (
	"bufio"
	environment "evie/env"
	"evie/values"
	"fmt"
	"os"
	"strconv"

	"time"

	"github.com/sanity-io/litter"
)

func SetupEnvironment(env *environment.Environment) {

	env.Variables["input"] = values.NativeFunctionValue{Value: ReadUserInput}
	env.Variables["print"] = values.NativeFunctionValue{Value: PrintStdOut}
	env.Variables["number"] = values.NativeFunctionValue{Value: ToNumber}
	env.Variables["string"] = values.NativeFunctionValue{Value: ToString}
	env.Variables["isNothing"] = values.NativeFunctionValue{Value: IsNothing}

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

func IsNothing(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing argument to isNothing function"}
	}

	for _, arg := range args {
		if arg.GetType() != "NothingValue" {
			return values.BooleanValue{Value: false}
		}
	}

	return values.BooleanValue{Value: true}
}

func ToNumber(args []values.RuntimeValue) values.RuntimeValue {
	val := args[0]

	switch value := val.(type) {
	case values.StringValue:
		number, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return values.ErrorValue{Value: err.Error()}
		}
		return values.NumberValue{Value: number}
	case values.NumberValue:
		return value
	default:
		return values.NumberValue{Value: 0}
	}
}

func ToString(args []values.RuntimeValue) values.RuntimeValue {
	arg := args[0]

	switch value := arg.(type) {
	case values.StringValue:
		return value
	case values.NumberValue:
		return values.StringValue{Value: value.GetStr()}
	case values.BooleanValue:
		return values.StringValue{Value: value.GetStr()}
	case values.ArrayValue:
		str := ""

		for _, item := range value.Value {
			str += item.GetStr()
			str += ", "
		}

		return values.StringValue{Value: str}
	default:
		return values.ErrorValue{Value: "Invalid conversion to string with type " + arg.GetType()}
	}
}

func ReadUserInput(args []values.RuntimeValue) values.RuntimeValue {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() { // Lee una línea completa
		return values.StringValue{Value: ""}
	}
	line := scanner.Text()
	return values.StringValue{Value: line}
}

func PrintValues(args []values.RuntimeValue, verboseStrings bool) {

	for _, arg := range args {
		if arg.GetType() == "ArrayValue" {
			fmt.Print("[ ")
			for _, item := range arg.(*values.ArrayValue).Value {
				PrintValues([]values.RuntimeValue{item}, true)
				fmt.Print(", ")
			}
			fmt.Print("] ")
		} else if arg.GetType() == "StringValue" {
			if verboseStrings {
				fmt.Print("'" + arg.GetStr() + "'")
			} else {
				fmt.Print(arg.GetStr())
			}
		} else if arg.GetType() == "NumberValue" {
			fmt.Print(arg.GetNumber())
		} else if arg.GetType() == "BooleanValue" {
			fmt.Print(arg.GetBool())
		} else if arg.GetType() == "DictionaryValue" {
			fmt.Print("{ ")
			for key, value := range arg.(values.DictionaryValue).Value {
				PrintValues([]values.RuntimeValue{values.StringValue{Value: key}}, true)
				fmt.Print(": ")
				PrintValues([]values.RuntimeValue{value}, true)
				fmt.Print(", ")
			}
			fmt.Print("}")
		} else {
			fmt.Print(arg.GetStr())
		}
	}

}

func PrintStdOut(args []values.RuntimeValue) values.RuntimeValue {
	PrintValues(args, false)
	print("\n")

	return values.BooleanValue{Value: true}
}

func GetArguments(args []values.RuntimeValue) values.RuntimeValue {

	vals := values.ArrayValue{Value: []values.RuntimeValue{}}

	sysargs := os.Args

	for _, arg := range sysargs {
		vals.Value = append(vals.Value, values.StringValue{Value: arg})
	}

	return &vals

}
