package runtime

import (
	environment "evie/env"
	"evie/values"
	"fmt"
	"strconv"
)

func IsError(val values.RuntimeValue) bool {
	return (val != nil && val.GetType() == "ErrorValue")
}
func (e *Evaluator) Panic(msg string, line int, env *environment.Environment) values.ErrorValue {
	output := "\n >>> DONT PANIC, but something went wrong at line " + fmt.Sprint(line) + " at module " + env.ModuleName + ":\n\t" + msg + "\n"

	if len(e.CallStack) > 0 {

		output += "\n >>> Very detailful Backtrace :D\n"
		for _, call := range e.CallStack {
			output += "\t" + call + "\n"
		}
	}
	return values.ErrorValue{Value: output}
}

func (e *Evaluator) CallStackEntry(mod string, line int) {
	e.CallStack = append([]string{"module -> " + mod + ": " + strconv.Itoa(line)}, e.CallStack...)
}

func (e *Evaluator) CallStackExit() {
	e.CallStack = e.CallStack[1:]
}
