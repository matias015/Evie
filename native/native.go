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

	env.PushScope()

	env.DeclareVar("input", values.NativeFunctionValue{Value: ReadUserInput})
	env.DeclareVar("print", values.NativeFunctionValue{Value: PrintStdOut})
	env.DeclareVar("number", values.NativeFunctionValue{Value: ToNumber})
	env.DeclareVar("string", values.NativeFunctionValue{Value: ToString})
	env.DeclareVar("isNothing", values.NativeFunctionValue{Value: IsNothing})

	env.DeclareVar("time", values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		now := time.Now()
		milliseconds := now.UnixMilli()
		return values.NumberValue{Value: float64(milliseconds)}
	}})

	env.DeclareVar("litter", values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		litter.Dump(args[0])
		return values.BoolValue{Value: true}
	}})

	env.DeclareVar("panic", values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		return values.ErrorValue{Value: "Panic! -> " + args[0].(values.StringValue).Value}
	}})

	arguments := []values.RuntimeValue{}

	for _, arg := range os.Args {
		arguments = append(arguments, values.StringValue{Value: arg})
	}

	// env.DeclareVar("getArgs", values.RuntimeValue{Value: GetArguments})
}

func IsNothing(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing argument to isNothing function"}
	}

	for _, arg := range args {
		if arg.GetType() != values.NothingType {
			return values.BoolValue{Value: false}
		}
	}

	return values.BoolValue{Value: true}
}

func ToNumber(args []values.RuntimeValue) values.RuntimeValue {
	val := args[0]

	switch val.GetType() {
	case values.StringType:
		number, err := strconv.ParseFloat(val.(values.StringValue).Value, 64)
		if err != nil {
			return values.ErrorValue{Value: err.Error()}
		}
		return values.NumberValue{Value: number}
	case values.NumberType:
		return val
	default:
		return values.ErrorValue{Value: "Invalid conversion to number with type " + val.GetType().String()}
	}
}

func ToString(args []values.RuntimeValue) values.RuntimeValue {
	value := args[0]

	switch value.GetType() {
	case values.StringType:
		return value
	case values.NumberType:
		return values.StringValue{Value: "value.String()"}
	case values.BoolType:
		return values.StringValue{Value: "value.String()"}
	case values.ArrayType:
		str := ""

		for _, item := range value.(*values.ArrayValue).Value {
			str += item.(values.StringValue).Value
			str += ", "
		}

		return values.StringValue{Value: str}
	default:
		return values.StringValue{Value: "Invalid conversion to string with type " + value.GetType().String()}
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

		valType := arg.GetType()

		if valType == values.ArrayType {
			fmt.Print("[ ")
			arr := arg.(*values.ArrayValue)
			for _, item := range arr.Value {
				PrintValues([]values.RuntimeValue{item}, true)
				fmt.Print(", ")
			}
			fmt.Print("] ")
		} else if valType == values.StringType {
			if verboseStrings {
				fmt.Print("'" + arg.(values.StringValue).Value + "'")
			} else {
				fmt.Print(arg.(values.StringValue).Value)
			}
		} else if valType == values.NumberType {
			fmt.Print(arg.(values.NumberValue).Value)
		} else if valType == values.BoolType {
			fmt.Print(arg.(values.BoolValue).Value)
		} else if valType == values.DictionaryType {
			fmt.Print("{ ")
			for key, value := range arg.(*values.DictionaryValue).Value {
				PrintValues([]values.RuntimeValue{values.StringValue{Value: key}}, true)
				fmt.Print(": ")
				PrintValues([]values.RuntimeValue{value}, true)
				fmt.Print(", ")
			}
			fmt.Print("}")
		} else {
			fmt.Print(valType.String())
		}
	}

}

func PrintStdOut(args []values.RuntimeValue) values.RuntimeValue {
	PrintValues(args, false)
	print("\n")

	return values.BoolValue{Value: true}
}

func GetArguments(args []values.RuntimeValue) values.RuntimeValue {

	vals := values.ArrayValue{Value: []values.RuntimeValue{}}

	sysargs := os.Args

	for _, arg := range sysargs {
		vals.Value = append(vals.Value, values.StringValue{Value: arg})
	}

	return &vals
}
