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
func (a *DictionaryValue) GetProp(name string) (RuntimeValue, error) {
	if name == "add" {
		return NativeFunctionValue{Value: DictAdd(a)}, nil
	} else if name == "remove" {
		return NativeFunctionValue{Value: DictRemove(a)}, nil
	} else if name == "has" {
		return NativeFunctionValue{Value: DictHas(a)}, nil
	} else {
		return NothingValue{}, fmt.Errorf("property %s does not exists", name)
	}
}

func DictAdd(a *DictionaryValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) < 2 {
			return ErrorValue{ErrorType: RuntimeError, Value: "Missing arguments for add method"}
		}

		if args[0].GetType() != StringType {
			return ErrorValue{ErrorType: RuntimeError, Value: "First argument (key) for add method must be a string"}
		}

		key := args[0].(StringValue).Value
		value := args[1]
		a.Value[key] = value
		return value
	}
}

func DictRemove(a *DictionaryValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) < 1 {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing arguments for remove method"}
		}
		if args[0].GetType() != StringType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument (key) for remove method must be a string"}
		}
		key := args[0].(StringValue).Value

		if _, ok := a.Value[key]; !ok {
			return ErrorValue{ErrorType: RuntimeError, Value: "Key " + key + " does not exist in dictionary"}
		}

		delete(a.Value, key)
		return BoolValue{Value: true}
	}
}

func DictHas(a *DictionaryValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) < 1 {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing arguments for has method"}
		}
		if args[0].GetType() != StringType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument (key) for has method must be a string"}
		}
		key := args[0].(StringValue).Value
		_, ok := a.Value[key]
		return BoolValue{Value: ok}
	}
}
