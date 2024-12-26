package values

import (
	"fmt"
	"strconv"
	"strings"
)

/*
----------------------------------------------------------
--- StringValue
----------------------------------------------------------
*/
// type StringValue struct {
// 	Value   string
// 	Mutable bool
// }

// func (s StringValue) GetStr() string     { return s.Value }
// func (s StringValue) GetNumber() float64 { return 0 }
// func (s StringValue) GetBool() bool {
// 	if s.Value == "" {
// 		return false
// 	} else {
// 		return true
// 	}
// }
// func (s StringValue) GetType() string {
// 	return "StringValue"
// }

func GetStringProp(value string, name string) (RuntimeValue, error) {

	props := map[string]RuntimeValue{
		"is": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				for _, arg := range args {
					if arg.Value.(string) == value {
						return RuntimeValue{Type: BoolType, Value: true}
					}
				}
				return RuntimeValue{Type: BoolType, Value: false}
			},
		},
		"addPaddingLeft": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				char := args[0].String()
				length := int(args[1].Value.(float64))
				actualLength := len(value)
				output := value
				for i := 0; i < length-actualLength; i++ {
					output = char + output
				}
				return RuntimeValue{Type: StringType, Value: output}
			},
		},
		"addPaddingRight": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				char := args[0].String()
				length := int(args[1].Value.(float64))
				actualLength := len(value)
				output := value
				for i := 0; i < length-actualLength; i++ {
					output = output + char
				}
				return RuntimeValue{Type: StringType, Value: output}
			},
		},
		"trim": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {

				needed := " "
				if len(args) > 0 {
					needed = args[0].String()
				}
				return RuntimeValue{Type: StringType, Value: strings.Trim(name, needed)}
			},
		},
		"toArray": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {

				sep := ""

				if len(args) > 0 {
					sep = args[0].String()
				}

				arr := ArrayValue{Value: make([]RuntimeValue, 0)}

				values := strings.Split(value, sep)

				for _, value := range values {
					arr.Value = append(arr.Value, RuntimeValue{Type: StringType, Value: value})
				}

				return RuntimeValue{Type: ArrayType, Value: &arr}
			},
		},
		"len": {Type: NativeFunctionType, Value: func(args []RuntimeValue) RuntimeValue {
			return RuntimeValue{Type: NumberType, Value: float64(len(value))}
		}},
		"slice": {Type: NativeFunctionType, Value: func(args []RuntimeValue) RuntimeValue {
			length := len(value)
			if len(args) == 2 {
				from := int(args[0].Value.(float64))
				to := int(args[1].Value.(float64))
				if to < 0 {
					to = length + to
				}
				if from < 0 {
					from = length + from
				}
				if from > length || to > length {
					return RuntimeValue{Type: ErrorType, Value: "Index out of range [" + strconv.FormatFloat(args[0].Value.(float64), 'f', -1, 64) + ":" + strconv.FormatFloat(args[1].Value.(float64), 'f', -1, 64) + "]"}
				}
				return RuntimeValue{Type: StringType, Value: value[from:to]}
			} else if len(args) == 1 {
				from := int(args[0].Value.(float64))
				if from < 0 {
					from = length + from
				}
				if from > length {
					return RuntimeValue{Type: ErrorType, Value: "Index out of range [" + strconv.FormatFloat(args[0].Value.(float64), 'f', -1, 64) + "]"}
				}
				return RuntimeValue{Type: StringType, Value: value[from:]}
			} else {
				return RuntimeValue{Type: StringType, Value: ""}
			}
		},
		},
	}

	p, ex := props[name]
	if !ex {
		return RuntimeValue{}, fmt.Errorf("Property %s not found", name)
	}

	return p, nil
}
