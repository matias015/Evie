package evruntime

import (
	environment "evie/env"
	"evie/values"
	"fmt"
)

func IsError(val values.RuntimeValue) bool {
	return (val.Type == values.ErrorType)
}
func (e *Evaluator) Panic(msg string, line int, env *environment.Environment) values.RuntimeValue {
	output := "\n >>> DONT PANIC, but something went wrong at line " + fmt.Sprint(line) + " at module " + env.ModuleName + ":\n\t" + msg + "\n"

	// if len(e.CallStack) > 0 {

	// 	output += "\n >>> Very detailful Backtrace :D\n"
	// 	for _, call := range e.CallStack {
	// 		output += "\t" + call + "\n"
	// 	}
	// }
	return values.RuntimeValue{Type: values.ErrorType, Value: output}
}
