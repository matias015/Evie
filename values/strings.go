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
type StringValue struct {
	Value string
}

func (a StringValue) GetString() string {
	return a.Value
}

func (a StringValue) GetType() ValueType {
	return StringType
}

func (s StringValue) GetProp(value *RuntimeValue, name string) (RuntimeValue, error) {

	props := map[string]RuntimeValue{
		"is": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				for _, arg := range args {
					if arg.(StringValue).Value == s.Value {
						return BoolValue{Value: true}
					}
				}
				return BoolValue{Value: false}
			},
		},
		"addPaddingLeft": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				char := args[0].(StringValue).Value
				length := int(args[1].(NumberValue).Value)
				actualLength := len(s.Value)
				output := s.Value
				for i := 0; i < length-actualLength; i++ {
					output = char + output
				}
				return StringValue{Value: output}
			},
		},
		"addPaddingRight": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				char := args[0].(StringValue).Value
				length := int(args[1].(NumberValue).Value)
				actualLength := len(s.Value)
				output := s.Value
				for i := 0; i < length-actualLength; i++ {
					output = output + char
				}
				return StringValue{Value: output}
			},
		},
		"trim": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {

				needed := " "
				if len(args) > 0 {
					needed = args[0].(StringValue).Value
				}
				return StringValue{Value: strings.Trim(s.Value, needed)}
			},
		},
		"toArray": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {

				sep := ""

				if len(args) > 0 {
					sep = args[0].(StringValue).Value
				}

				arr := ArrayValue{Value: make([]RuntimeValue, 0)}

				values := strings.Split(s.Value, sep)

				for _, value := range values {
					arr.Value = append(arr.Value, StringValue{Value: value})
				}

				return &arr
			},
		},
		"len": NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
			return NumberValue{Value: float64(len(s.Value))}
		}},
		"slice": NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
			length := len(s.Value)
			if len(args) == 2 {
				from := int(args[0].(NumberValue).Value)
				to := int(args[1].(NumberValue).Value)
				if to < 0 {
					to = length + to
				}
				if from < 0 {
					from = length + from
				}
				if from > length || to > length {
					return ErrorValue{Value: "Index out of range [" + strconv.Itoa(from) + ":" + strconv.Itoa(to) + "]"}
				}
				return ErrorValue{Value: s.Value[from:to]}
			} else if len(args) == 1 {
				from := int(args[0].(NumberValue).Value)
				if from < 0 {
					from = length + from
				}
				if from > length {
					return ErrorValue{Value: "Index out of range [" + strconv.Itoa(from) + "]"}
				}
				return StringValue{Value: s.Value[from:]}
			} else {
				return StringValue{Value: ""}
			}
		},
		},
	}

	p, ex := props[name]
	if !ex {
		return NothingValue{}, fmt.Errorf("Property %s not found", name)
	}

	return p, nil
}
