package fsLib

import (
	"bufio"
	environment "evie/env"
	"evie/values"
	"log"
	"os"
)

type FileValue struct {
	Value   *os.File
	Scanner *bufio.Scanner
	Closed  bool
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

func (s *FileValue) GetProp(v *values.RuntimeValue, name string) values.RuntimeValue {
	props := map[string]values.RuntimeValue{}

	props = map[string]values.RuntimeValue{
		"readLine": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {

				if s.Closed {
					return values.ErrorValue{Value: "File is closed"}
				}

				success := s.Scanner.Scan()
				if !success {
					if err := s.Scanner.Err(); err != nil {
						return values.ErrorValue{Value: err.Error()}
					} else {
						return values.StringValue{Value: "EOF"}
					}
				}
				return values.StringValue{Value: s.Scanner.Text()}
			},
		},
		"appendLine": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				_, err := s.Value.WriteString(args[0].GetStr() + "\n")

				if err != nil {
					return values.ErrorValue{Value: err.Error()}
				}
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
				s.Closed = true
				s.Value.Close()
				return values.BooleanValue{Value: true}
			},
		},
		"seek": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				arg := args[0]
				_, err := s.Value.Seek(0, int(arg.GetNumber()))

				if err != nil {
					return values.ErrorValue{Value: err.Error()}
				}

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

	env.DeclareVar("File", GetFileStruct())

	structValue := values.StructValue{}
	structValue.Properties = []string{"fs"}

	structValue.Methods = make(map[string]values.RuntimeValue)

	structValue.Methods["read"] = values.NativeFunctionValue{Value: ReadFile}
	structValue.Methods["open"] = values.NativeFunctionValue{Value: OpenFile}
	structValue.Methods["exists"] = values.NativeFunctionValue{Value: FileExists}
	structValue.Methods["remove"] = values.NativeFunctionValue{Value: RemoveFile}

	obj := values.ObjectValue{}
	obj.Struct = structValue
	obj.Value = make(map[string]values.RuntimeValue)

	env.DeclareVar("fs", obj)

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

	file, err := os.OpenFile(args[0].GetStr(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}

	f := FileValue{}

	f.Value = file
	f.Scanner = bufio.NewScanner(file)
	f.Closed = false

	return &f
}

func FileExists(args []values.RuntimeValue) values.RuntimeValue {
	arg := args[0]
	_, err := os.Stat(arg.GetStr())
	if err != nil {
		return values.BooleanValue{Value: false}
	}
	return values.BooleanValue{Value: true}
}

func RemoveFile(args []values.RuntimeValue) values.RuntimeValue {
	err := os.Remove(args[0].GetStr())
	if err != nil {
		return values.ErrorValue{Value: err.Error()}
	}
	return values.BooleanValue{Value: true}
}
