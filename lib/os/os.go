package osLib

import (
	environment "evie/env"
	"evie/values"
	"os"
	"os/exec"
	"path/filepath"
)

func Load(env *environment.Environment) {

	ns := values.NamespaceValue{Value: make(map[string]values.RuntimeValue)}

	ns.Value["exec"] = values.NativeFunctionValue{Value: Exec}

	ns.Value["changeDir"] = values.NativeFunctionValue{Value: ChangeWD}
	ns.Value["getWDir"] = values.NativeFunctionValue{Value: GetWD}

	ns.Value["getProcessId"] = values.NativeFunctionValue{Value: GetProcessId}
	ns.Value["processExists"] = values.NativeFunctionValue{Value: GetProcessById}

	ns.Value["kill"] = values.NativeFunctionValue{Value: KillProcess}

	ns.Value["PATH_SEPARATOR"] = values.StringValue{Value: string(filepath.Separator)}

	env.ForceDeclare("os", ns)

}

func Exec(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) == 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Exec method needs at least 1 argument"}
	}

	strArgs := make([]string, len(args))

	for i, arg := range args {
		if arg.GetType() == values.StringType || arg.GetType() == values.NumberType {
			strArgs[i] = arg.GetString()
		}
	}

	cmd := exec.Command(strArgs[0], strArgs[1:]...)

	output, err := cmd.Output()

	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}

	return values.StringValue{Value: string(output)}
}

func GetProcessId(args []values.RuntimeValue) values.RuntimeValue {
	return values.NumberValue{Value: float64(os.Getpid())}
}

func ChangeWD(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) < 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file create method"}
	} else {
		if args[0].GetType() != values.StringType {
			return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file create method"}
		}
	}

	dir := args[0].(values.StringValue).Value

	err := os.Chdir(dir)

	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}

	return values.NothingValue{}

}

func GetWD(args []values.RuntimeValue) values.RuntimeValue {

	dir, err := os.Getwd()

	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}

	return values.StringValue{Value: dir}

}

func KillProcess(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) < 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for process kill method"}
	} else {
		if args[0].GetType() != values.NumberType {
			return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Invalid arguments for process kill method, process id should be a number"}
		}
	}

	pid := int(args[0].GetNumber())

	pc, err := os.FindProcess(pid)

	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	} else {
		err := pc.Kill()

		if err != nil {
			return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
		} else {
			return values.BoolValue{Value: true}
		}
	}

}

func GetProcessById(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) < 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file Get process by id method"}
	} else {
		if args[0].GetType() != values.NumberType {
			return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file Get process by id method"}
		}
	}

	pid := args[0].GetNumber()

	_, err := os.FindProcess(int(pid))

	if err != nil {
		return values.BoolValue{Value: false}
	} else {
		return values.BoolValue{Value: true}
	}

}
