package native

import (
	environment "evie/env"
	"evie/values"
	"fmt"
	"os"

	"time"

	"github.com/sanity-io/litter"
)

func SetupEnvironment(env *environment.Environment) {

	env.PushScope()

	// env.DeclareVar("input", values.RuntimeValue{Type: values.NativeFunctionType, Value: ReadUserInput})
	env.DeclareVar("print", values.RuntimeValue{Type: values.NativeFunctionType, Value: PrintStdOut})
	// env.DeclareVar("number", values.RuntimeValue{Type: values.NativeFunctionType, Value: ToNumber})
	// env.DeclareVar("string", values.RuntimeValue{Type: values.NativeFunctionType, Value: ToString})
	env.DeclareVar("isNothing", values.RuntimeValue{Type: values.NativeFunctionType, Value: IsNothing})

	env.DeclareVar("time", values.RuntimeValue{Type: values.NativeFunctionType, Value: func(args []values.RuntimeValue) values.RuntimeValue {
		now := time.Now()
		milliseconds := now.UnixMilli()
		return values.RuntimeValue{Type: values.NumberType, Value: float64(milliseconds)}
	}})

	env.DeclareVar("litter", values.RuntimeValue{Type: values.NativeFunctionType, Value: func(args []values.RuntimeValue) values.RuntimeValue {
		litter.Dump(args[0])
		return values.RuntimeValue{Type: values.BoolType, Value: true}
	}})

	env.DeclareVar("panic", values.RuntimeValue{Type: values.NativeFunctionType, Value: func(args []values.RuntimeValue) values.RuntimeValue {
		return values.RuntimeValue{Type: values.ErrorType, Value: "Panic! -> " + args[0].Value.(string)}
	}})

	arguments := []values.RuntimeValue{}

	for _, arg := range os.Args {
		arguments = append(arguments, values.RuntimeValue{Type: values.StringType, Value: arg})
	}

	// env.DeclareVar("getArgs", values.RuntimeValue{Type: values.NativeFunctionType, Value: GetArguments})
}

func IsNothing(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.RuntimeValue{Type: values.ErrorType, Value: "Missing argument to isNothing function"}
	}

	for _, arg := range args {
		if arg.Type != values.NothingType {
			return values.RuntimeValue{Type: values.BoolType, Value: false}
		}
	}

	return values.RuntimeValue{Type: values.BoolType, Value: true}
}

// func ToNumber(args []values.RuntimeValue) values.RuntimeValue {
// 	val := args[0]

// 	switch value := val.(type) {
// 	case values.StringValue:
// 		number, err := strconv.ParseFloat(value.Value, 64)
// 		if err != nil {
// 			return values.ErrorValue{Value: err.Error()}
// 		}
// 		return values.NumberValue{Value: number}
// 	case values.NumberValue:
// 		return value
// 	default:
// 		return values.NumberValue{Value: 0}
// 	}
// }

// func ToString(args []values.RuntimeValue) values.RuntimeValue {
// 	arg := args[0]

// 	switch value := arg.(type) {
// 	case values.StringValue:
// 		return value
// 	case values.NumberValue:
// 		return values.StringValue{Value: value.GetStr()}
// 	case values.BooleanValue:
// 		return values.StringValue{Value: value.GetStr()}
// 	case values.ArrayValue:
// 		str := ""

// 		for _, item := range value.Value {
// 			str += item.GetStr()
// 			str += ", "
// 		}

// 		return values.StringValue{Value: str}
// 	default:
// 		return values.ErrorValue{Value: "Invalid conversion to string with type " + arg.GetType()}
// 	}
// }

// func ReadUserInput(args []values.RuntimeValue) values.RuntimeValue {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	if !scanner.Scan() { // Lee una línea completa
// 		return values.StringValue{Value: ""}
// 	}
// 	line := scanner.Text()
// 	return values.StringValue{Value: line}
// }

func PrintValues(args []values.RuntimeValue, verboseStrings bool) {

	for _, arg := range args {
		if arg.Type == values.ArrayType {
			fmt.Print("[ ")
			arr := arg.Value.(*values.ArrayValue)
			for _, item := range arr.Value {
				PrintValues([]values.RuntimeValue{item}, true)
				fmt.Print(", ")
			}
			fmt.Print("] ")
		} else if arg.Type == values.StringType {
			if verboseStrings {
				fmt.Print("'" + arg.Value.(string) + "'")
			} else {
				fmt.Print(arg.Value.(string))
			}
		} else if arg.Type == values.NumberType {
			fmt.Print(arg.Value.(float64))
		} else if arg.Type == values.BoolType {
			fmt.Print(arg.Value.(bool))
		} else if arg.Type == values.DictionaryType {
			fmt.Print("{ ")
			for key, value := range arg.Value.(*values.DictionaryValue).Value {
				PrintValues([]values.RuntimeValue{{Type: values.StringType, Value: key}}, true)
				fmt.Print(": ")
				PrintValues([]values.RuntimeValue{value}, true)
				fmt.Print(", ")
			}
			fmt.Print("}")
		} else {
			fmt.Print(arg.Value)
		}
	}

}

func PrintStdOut(args []values.RuntimeValue) values.RuntimeValue {
	PrintValues(args, false)
	print("\n")

	return values.RuntimeValue{Type: values.BoolType, Value: true}
}

// func GetArguments(args []values.RuntimeValue) values.RuntimeValue {

// 	vals := values.ArrayValue{Value: []values.RuntimeValue{}}

// 	sysargs := os.Args

// 	for _, arg := range sysargs {
// 		vals.Value = append(vals.Value, values.StringValue{Value: arg})
// 	}

// 	return &vals

// }
