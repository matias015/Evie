package osLib

import (
	environment "evie/env"
	"evie/values"
	"os"
	"os/exec"
)

func Load(env *environment.Environment) {

	ns := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}

	ns.Value["exec"] = values.NativeFunctionValue{Value: Exec}
	ns.Value["getProcessId"] = values.NativeFunctionValue{Value: GetProcessId}
	ns.Value["createDir"] = values.NativeFunctionValue{Value: CreateDir}

	env.DeclareVar("os", ns)

}

func Exec(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing arguments"}
	}

	strArgs := make([]string, len(args))

	for i, arg := range args {
		strArgs[i] = arg.(values.StringValue).Value
	}

	cmd := exec.Command(strArgs[0], strArgs[1:]...)

	output, err := cmd.Output()

	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}

	return values.StringValue{Value: string(output)}
}

func GetProcessId(args []values.RuntimeValue) values.RuntimeValue {
	return values.NumberValue{Value: float64(os.Getpid())}
}

func CreateDir(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) == 0 {
		return values.ErrorValue{Value: "Missing arguments"}
	}

	path := args[0].(values.StringValue).Value

	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}

	return values.BoolValue{Value: true}
}
