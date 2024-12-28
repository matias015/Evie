package values

import (
	"fmt"
)

type DictionaryValue struct {
	Value map[string]RuntimeValue
}

func (a DictionaryValue) GetString() string {
	return "Dictionary"
}

func (a DictionaryValue) GetBool() bool {
	return true
}
func (a DictionaryValue) GetNumber() float64 {
	return 0
}
func (a DictionaryValue) GetType() ValueType {
	return DictionaryType
}
func (a *DictionaryValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	props := map[string]RuntimeValue{

		"add": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				if len(args) < 2 {
					return ErrorValue{Value: "Missing arguments for add method"}
				}

				if args[0].GetType() != StringType {
					return ErrorValue{Value: "First argument (key) for add method must be a string"}
				}

				key := args[0].(StringValue).Value
				value := args[1]
				a.Value[key] = value
				return value
			},
		},
		"remove": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				if len(args) < 1 {
					return ErrorValue{Value: "Missing arguments for remove method"}
				}
				if args[0].GetType() != StringType {
					return ErrorValue{Value: "First argument (key) for remove method must be a string"}
				}
				key := args[0].(StringValue).Value

				if _, ok := a.Value[key]; !ok {
					return ErrorValue{Value: "Key " + key + " does not exist in dictionary"}
				}

				delete(a.Value, key)
				return BoolValue{Value: true}
			},
		},
		"has": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				if len(args) < 1 {
					return ErrorValue{Value: "Missing arguments for has method"}
				}
				if args[0].GetType() != StringType {
					return ErrorValue{Value: "First argument (key) for has method must be a string"}
				}
				key := args[0].(StringValue).Value
				_, ok := a.Value[key]
				return BoolValue{Value: ok}
			},
		},
	}

	prop, ex := props[name]

	if !ex {
		return NothingValue{}, fmt.Errorf("property %s does not exists", name)
	}

	return prop, nil
}
