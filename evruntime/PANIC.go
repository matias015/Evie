package evruntime

import (
	environment "evie/env"
	"evie/values"
	"fmt"
)

func IsError(val values.RuntimeValue) bool {
	return (val.GetType() == values.ErrorType)
}
func (e *Evaluator) Panic(msg string, line int, env *environment.Environment) values.RuntimeValue {
	output := "\n >>> DONT PANIC, but something went wrong at line " + fmt.Sprint(line) + " at module " + env.ModuleName + ":\n\t" + msg + "\n"

	if len(e.CallStack) > 0 {
		output += "\n >>> Detailed callstack:\n"
		for i := len(e.CallStack) - 1; i >= 0; i-- {
			output += fmt.Sprintf("\t%s\n", e.CallStack[i])
		}
	}

	return values.ErrorValue{Value: output}
}
