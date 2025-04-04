package fsLib

import (
	"bufio"
	environment "evie/env"
	"evie/values"
	"os"
)

func GetFileStruct() values.StructValue {
	fileStruct := values.StructValue{}
	fileStruct.Properties = []string{"descriptor"}
	fileStruct.Methods = make(map[string]values.RuntimeValue)
	return fileStruct
}

func Load(env *environment.Environment) {

	ns := values.NamespaceValue{}

	methods := make(map[string]values.RuntimeValue, 4)

	methods["read"] = values.NativeFunctionValue{Value: ReadFile}
	methods["open"] = values.NativeFunctionValue{Value: OpenFile}
	methods["exists"] = values.NativeFunctionValue{Value: FileExists}
	methods["remove"] = values.NativeFunctionValue{Value: RemoveFile}
	methods["create"] = values.NativeFunctionValue{Value: CreateFile}

	methods["createDir"] = values.NativeFunctionValue{Value: CreateDir}
	methods["existsDir"] = values.NativeFunctionValue{Value: FileExistsDir}
	methods["removeDir"] = values.NativeFunctionValue{Value: RemoveDir}
	methods["moveOrRenameDir"] = values.NativeFunctionValue{Value: MoveDir}
	ns.Value = methods

	env.ForceDeclare("fs", ns)

}

func ReadFile(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) < 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file read method"}
	} else {
		if args[0].GetType() != values.StringType {
			return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for file read method should be string"}
		}
	}

	bytes, err := os.ReadFile(args[0].(values.StringValue).Value)
	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}

	// imprimir el string
	return values.StringValue{Value: string(bytes)}
}

func OpenFile(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) < 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file open method"}
	} else {
		if args[0].GetType() != values.StringType {
			return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for file open method should be string"}
		}
	}

	file, err := os.OpenFile(args[0].(values.StringValue).Value, os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}

	f := FileValue{}

	f.Value = file
	f.Scanner = bufio.NewScanner(file)
	f.Closed = false

	return f
}

func FileExists(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) < 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file exists method"}
	} else {
		if args[0].GetType() != values.StringType {
			return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for file exists method should be string"}
		}
	}

	arg := args[0]

	_, err := os.Stat(arg.(values.StringValue).Value)

	if err != nil {
		return values.BoolValue{Value: false}
	}

	return values.BoolValue{Value: true}
}

func RemoveFile(args []values.RuntimeValue) values.RuntimeValue {
	if len(args) < 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file remove method"}
	} else {
		if args[0].GetType() != values.StringType {
			return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for file remove method should be string"}
		}
	}
	err := os.Remove(args[0].(values.StringValue).Value)
	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}
	return values.BoolValue{Value: true}
}

func CreateFile(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) == 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for file create method"}
	}

	if args[0].GetType() != values.StringType {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for file create method must be a string"}
	}

	file, err := os.Create(args[0].(values.StringValue).Value)
	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}
	return &FileValue{Value: file, Closed: false, Scanner: bufio.NewScanner(file)}
}

// Directory

func CreateDir(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) == 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for directory create method"}
	}

	if args[0].GetType() != values.StringType {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for directory create method must be a string"}
	}

	err := os.Mkdir(args[0].(values.StringValue).Value, 0755)
	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}
	return values.BoolValue{Value: true}
}

func RemoveDir(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) == 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for directory remove method"}
	}

	if args[0].GetType() != values.StringType {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for directory remove method must be a string"}
	}

	err := os.RemoveAll(args[0].(values.StringValue).Value)

	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}
	return values.BoolValue{Value: true}
}

func FileExistsDir(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) == 0 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for directory exists method"}
	}

	if args[0].GetType() != values.StringType {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "First argument for directory exists method must be a string"}
	}

	_, err := os.Stat(args[0].(values.StringValue).Value)
	if err != nil {
		return values.BoolValue{Value: false}
	}

	return values.BoolValue{Value: true}
}

func MoveDir(args []values.RuntimeValue) values.RuntimeValue {

	if len(args) < 1 {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Missing arguments for directory move method"}
	}

	if args[0].GetType() != values.StringType || args[1].GetType() != values.StringType {
		return values.ErrorValue{ErrorType: values.InvalidArgumentError, Value: "Arguments for directory move method must be a string"}
	}

	err := os.Rename(args[0].(values.StringValue).Value, args[1].(values.StringValue).Value)
	if err != nil {
		return values.ErrorValue{ErrorType: values.RuntimeError, Value: err.Error()}
	}
	return values.BoolValue{Value: true}

}
