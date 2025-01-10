package fsLib

import (
	"bufio"
	"evie/values"
	"os"
)

type FileValue struct {
	Value   *os.File
	Scanner *bufio.Scanner
	Closed  bool
}

func (s FileValue) GetType() values.ValueType {
	return values.FileType
}
func (a FileValue) GetNumber() float64 {
	return 1
}
func (a FileValue) GetBool() bool {
	return !a.Closed
}

func (s FileValue) GetString() string {
	return "File Value"
}

func (s FileValue) GetProp(name string) (values.RuntimeValue, error) {
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
				_, err := s.Value.WriteString(args[0].(values.StringValue).Value + "\n")

				if err != nil {
					return values.ErrorValue{Value: err.Error()}
				}
				return values.BoolValue{Value: true}
			},
		},
		"append": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				s.Value.WriteString(args[0].(values.StringValue).Value)
				return values.BoolValue{Value: true}
			},
		},
		"close": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				s.Closed = true
				s.Value.Close()
				return values.BoolValue{Value: true}
			},
		},
		"seek": values.NativeFunctionValue{
			Value: func(args []values.RuntimeValue) values.RuntimeValue {
				arg := args[0]
				_, err := s.Value.Seek(0, int(arg.(values.NumberValue).Value))

				if err != nil {
					return values.ErrorValue{Value: err.Error()}
				}

				return values.BoolValue{Value: true}
			},
		},
	}

	p, ex := props[name]

	if !ex {
		return values.ErrorValue{Value: "property " + name + " does not exists"}, nil
	}

	return p, nil
}
