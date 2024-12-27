package values

import "fmt"

type DictionaryValue struct {
	Value map[string]RuntimeValue
}

func (a DictionaryValue) GetString() string {
	return "Dictionary"
}

func (a DictionaryValue) GetType() ValueType {
	return DictionaryType
}
func (a DictionaryValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	props := map[string]RuntimeValue{

		"add": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].(StringValue).Value
				value := args[1]
				a.Value[key] = value
				return BoolValue{Value: true}
			},
		},
		"remove": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].(StringValue).Value
				delete(a.Value, key)
				return BoolValue{Value: true}
			},
		},
		"has": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
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
