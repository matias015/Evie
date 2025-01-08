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

	env.ForceDeclare("RuntimeError", values.StringValue{Value: "RuntimeError"})
	env.ForceDeclare("TypeError", values.StringValue{Value: "TypeError"})
	env.ForceDeclare("InvalidIndexError", values.StringValue{Value: "InvalidIndexError"})
	env.ForceDeclare("IdentifierError", values.StringValue{Value: "IdentifierError"})
	env.ForceDeclare("ZeroDivisionError", values.StringValue{Value: "ZeroDivisionError"})
	env.ForceDeclare("InvalidArgumentError", values.StringValue{Value: "InvalidArgumentError"})
	env.ForceDeclare("InvalidConversionError", values.StringValue{Value: "InvalidConversionError"})
	env.ForceDeclare("CircularImportError", values.StringValue{Value: "CircularImportError"})
	env.ForceDeclare("PropertyError", values.StringValue{Value: "PropertyError"})

	env.DeclareVar("ErrorObject", values.StructValue{
		Name:    "ErrorObject",
		Methods: make(map[string]values.RuntimeValue),
		Properties: []string{
			"message", "type",
		},
	})

	env.DeclareVar("input", values.NativeFunctionValue{Value: ReadUserInput})
	env.DeclareVar("print", values.NativeFunctionValue{Value: PrintStdOut})
	env.DeclareVar("number", values.NativeFunctionValue{Value: ToNumber})
	env.DeclareVar("int", values.NativeFunctionValue{Value: ToInteger})
	env.DeclareVar("string", values.NativeFunctionValue{Value: ToString})
	env.DeclareVar("bool", values.NativeFunctionValue{Value: ToBool})
	env.DeclareVar("isNothing", values.NativeFunctionValue{Value: IsNothing})
	env.DeclareVar("type", values.NativeFunctionValue{Value: Type})

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
		if len(args) == 0 {
			return values.ErrorValue{Value: "Why?"}
		}
		if args[0].GetType() == values.StringType {
			return values.ErrorValue{Value: "Panic! -> " + args[0].GetString()}
		} else if args[0].GetType() == values.ObjectType {
			obj := args[0].(*values.ObjectValue)
			return values.ErrorValue{Value: obj.Value["message"].GetString(), ErrorType: obj.Value["type"].GetString()}
		} else {
			return values.ErrorValue{Value: "Panic while trying to panic, because panic argument is not valid"}
		}
	}})

	arguments := values.ArrayValue{Value: []values.RuntimeValue{}}

	for _, arg := range os.Args {
		arguments.Value = append(arguments.Value, values.StringValue{Value: arg})
	}

	env.DeclareVar("getArgs", values.NativeFunctionValue{Value: func(args []values.RuntimeValue) values.RuntimeValue {
		return &arguments
	}})

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
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing argument to number parse function"}
	}

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
func ToInteger(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing argument to integer parse function"}
	}

	val := args[0]

	switch val.GetType() {
	case values.StringType:
		number, err := strconv.Atoi(val.(values.StringValue).Value)
		if err != nil {
			return values.ErrorValue{Value: err.Error()}
		}
		return values.NumberValue{Value: float64(number)}
	case values.NumberType:
		return values.NumberValue{Value: float64(int(val.(values.NumberValue).Value))}
	default:
		return values.ErrorValue{Value: "Invalid conversion to number with type " + val.GetType().String()}
	}
}

func ToString(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing argument to string parse function"}
	}
	value := args[0]

	return values.StringValue{Value: value.GetString()}
}

func ToBool(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing argument to bool parse function"}
	}
	value := args[0]

	return values.BoolValue{Value: value.GetBool()}
}

func ReadUserInput(args []values.RuntimeValue) values.RuntimeValue {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
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
		} else if valType == values.ErrorType {
			fmt.Print("'" + arg.(values.ErrorValue).Value + "'")
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

func Type(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing argument for type function"}
	}
	if args[0].GetType() == values.ObjectType {
		return values.StringValue{Value: args[0].(*values.ObjectValue).Struct.Name}
	}
	return values.StringValue{Value: args[0].GetType().String()}
}
