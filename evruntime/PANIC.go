package evruntime

import (
	environment "evie/env"
	"evie/values"
	"fmt"
)

func (e Evaluator) Panic(errorType string, msg string, line int, env *environment.Environment) values.ErrorValue {

	errorStruct, _ := env.GetVar("ErrorObject")

	errProperties := make(map[string]values.RuntimeValue)

	errProperties["message"] = values.StringValue{Value: msg}
	errProperties["type"] = values.StringValue{Value: errorType}
	errProperties["line"] = values.NumberValue{Value: float64(line)}
	errProperties["module"] = values.StringValue{Value: env.ModuleName}

	callStack := &values.ArrayValue{Value: make([]values.RuntimeValue, 0)}

	for i := len(e.CallStack.Items) - 1; i >= 0; i-- {
		callStack.Value = append(callStack.Value, values.StringValue{Value: e.CallStack.Items[i].String()})
	}

	errProperties["callstack"] = callStack

	return values.ErrorValue{
		Object: &values.ObjectValue{
			Struct: errorStruct.(values.StructValue),
			Value:  errProperties,
		},
	}
}

func (e Evaluator) PrintError(err values.ErrorValue) {

	errValue := err.Object

	output := "\n >>> DONT PANIC, but something went wrong at line " + errValue.Value["line"].GetString() + " at module " + errValue.Value["module"].GetString() + ":\n\t " + errValue.Value["type"].GetString() + ": " + errValue.Value["message"].GetString() + "\n"

	cstack := errValue.Value["callstack"].(*values.ArrayValue).Value

	if len(cstack) > 0 {
		output += "\n >>> Detailed callstack:\n"
		for _, item := range cstack {
			output += fmt.Sprintf("\t%s\n", item.GetString())
		}
	}

	fmt.Println(output)
}
