package fsLib

import (
	"bufio"
	environment "evie/env"
	"evie/values"
	"log"
	"os"
)

func GetFileStruct() values.StructValue {
	fileStruct := values.StructValue{}
	fileStruct.Properties = []string{"descriptor"}
	fileStruct.Methods = make(map[string]values.RuntimeValue)
	return fileStruct
}

func Load(env *environment.Environment) {

	env.DeclareVar("File", GetFileStruct())

	ns := values.NamespaceValue{}

	methods := make(map[string]values.RuntimeValue, 4)

	methods["read"] = values.NativeFunctionValue{Value: ReadFile}
	methods["open"] = values.NativeFunctionValue{Value: OpenFile}
	methods["exists"] = values.NativeFunctionValue{Value: FileExists}
	methods["remove"] = values.NativeFunctionValue{Value: RemoveFile}

	ns.Value = methods

	env.DeclareVar("fs", ns)

}

func ReadFile(args []values.RuntimeValue) values.RuntimeValue {
	datosComoBytes, err := os.ReadFile(args[0].(values.StringValue).Value)
	if err != nil {
		log.Fatal(err)
	}
	// convertir el arreglo a string
	datosComoString := string(datosComoBytes)
	// imprimir el string
	return values.StringValue{Value: datosComoString}
}

func OpenFile(args []values.RuntimeValue) values.RuntimeValue {

	file, err := os.OpenFile(args[0].(values.StringValue).Value, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}

	f := FileValue{}

	f.Value = file
	f.Scanner = bufio.NewScanner(file)
	f.Closed = false

	return f
}

func FileExists(args []values.RuntimeValue) values.RuntimeValue {
	arg := args[0]
	_, err := os.Stat(arg.(values.StringValue).Value)
	if err != nil {
		return values.BoolValue{Value: false}
	}
	return values.BoolValue{Value: true}
}

func RemoveFile(args []values.RuntimeValue) values.RuntimeValue {
	err := os.Remove(args[0].(values.StringValue).Value)
	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}
	return values.BoolValue{Value: true}
}
