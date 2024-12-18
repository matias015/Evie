package fsLib

import (
	"bufio"
	environment "evie/env"
	"evie/values"
	"log"
	"os"
)

type FileValue struct {
	Value  *os.File
	Reader *bufio.Reader
}

func (s FileValue) SetValue(value values.RuntimeValue) {}

func (s FileValue) GetStr() string     { return "" }
func (s FileValue) GetNumber() float64 { return 1 }
func (s FileValue) GetBool() bool {
	return false
}
func (s FileValue) GetType() string {
	return "FileValue"
}

func (s FileValue) GetProp(name string) values.RuntimeValue {
	props := map[string]values.RuntimeValue{}

	props = map[string]values.RuntimeValue{
		"readLine": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {

				t, err := s.Reader.ReadString('\n')
				if err != nil {
					return values.StringValue{Value: err.Error()}
				}
				return values.StringValue{Value: t[:len(t)-1]}
			},
		},
		"appendLine": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				s.Value.WriteString(args[0].GetStr() + "\n")
				return values.BooleanValue{Value: true}
			},
		},
		"append": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				s.Value.WriteString(args[0].GetStr())
				return values.BooleanValue{Value: true}
			},
		},
		"close": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				s.Value.Close()
				return values.BooleanValue{Value: true}
			},
		},
	}

	return props[name]
}

func GetFileStruct() values.StructValue {
	fileStruct := values.StructValue{}
	fileStruct.Properties = []string{"descriptor"}
	fileStruct.Methods = make(map[string]values.RuntimeValue)
	return fileStruct
}

func Load(env *environment.Environment) {

	env.Variables["File"] = GetFileStruct()

	structValue := values.StructValue{}
	structValue.Properties = []string{"fs"}

	structValue.Methods = make(map[string]values.RuntimeValue)

	structValue.Methods["readAll"] = values.NativeFunctionValue{Value: ReadFile}
	structValue.Methods["open"] = values.NativeFunctionValue{Value: OpenFile}
	obj := values.ObjectValue{}
	obj.Struct = structValue
	obj.Value = make(map[string]values.RuntimeValue)

	env.Variables["fs"] = obj

}

func ReadFile(args []values.RuntimeValue) values.RuntimeValue {
	datosComoBytes, err := os.ReadFile(args[0].GetStr())
	if err != nil {
		log.Fatal(err)
	}
	// convertir el arreglo a string
	datosComoString := string(datosComoBytes)
	// imprimir el string
	return values.StringValue{Value: datosComoString}
}

func OpenFile(args []values.RuntimeValue) values.RuntimeValue {

	file, err := os.OpenFile(args[0].GetStr(), os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}

	f := FileValue{}

	f.Value = file
	f.Reader = bufio.NewReader(file)

	return f
}
