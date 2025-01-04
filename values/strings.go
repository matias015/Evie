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

func (a StringValue) GetNumber() float64 {
	return 1
}
func (a StringValue) GetBool() bool {
	return a.Value != ""
}
func (a StringValue) GetString() string {
	return a.Value
}

func (a StringValue) GetType() ValueType {
	return StringType
}

func (s StringValue) GetProp(name string) (RuntimeValue, error) {

	if name == "len" {
		return NativeFunctionValue{Value: StringLen(s)}, nil
	} else if name == "is" {
		return NativeFunctionValue{Value: StringIs(s)}, nil
	} else if name == "addPaddingLeft" {
		return NativeFunctionValue{Value: StringPadLeft(s)}, nil
	} else if name == "addPaddingRight" {
		return NativeFunctionValue{Value: StringPadRight(s)}, nil
	} else if name == "toArray" {
		return NativeFunctionValue{Value: StringToArray(s)}, nil
	} else if name == "trim" {
		return NativeFunctionValue{Value: StringTrim(s)}, nil
	} else if name == "slice" {
		return NativeFunctionValue{Value: StringSlice(s)}, nil
	} else {
		return nil, fmt.Errorf("property %s does not exists", name)
	}
}

func StringIs(s StringValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		for _, arg := range args {
			if arg.(StringValue).Value == s.Value && arg.GetType() == StringType {
				return BoolValue{Value: true}
			}
		}
		return BoolValue{Value: false}
	}
}

func StringPadLeft(s StringValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) < 2 {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing arguments for addPaddingLeft method"}
		}

		if args[0].GetType() != StringType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument for addPaddingLeft method must be a string"}
		}

		if args[1].GetType() != NumberType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Second argument for addPaddingLeft method must be a number"}
		}

		char := args[0].GetString()
		length := int(args[1].(NumberValue).Value)
		output := s.Value
		for i := 0; i < length-len(s.Value); i++ {
			output = char + output
		}

		return StringValue{Value: output}
	}
}
func StringPadRight(s StringValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) < 2 {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing arguments for addPaddingLeft method"}
		}

		if args[0].GetType() != StringType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument for addPaddingLeft method must be a string"}
		}

		if args[1].GetType() != NumberType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Second argument for addPaddingLeft method must be a number"}
		}

		char := args[0].(StringValue).Value
		length := int(args[1].(NumberValue).Value)
		output := s.Value
		for i := 0; i < length-len(s.Value); i++ {
			output = output + char
		}
		return StringValue{Value: output}
	}
}

func StringTrim(s StringValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		needed := " "
		if len(args) > 0 {
			if args[0].GetType() != StringType {
				return ErrorValue{ErrorType: InvalidArgumentError, Value: "Argument for trim method must be a string"}
			}
			needed = args[0].(StringValue).Value
		}
		return StringValue{Value: strings.Trim(s.Value, needed)}
	}
}
func StringToArray(s StringValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {

		sep := ""

		if len(args) > 0 {
			if args[0].GetType() != StringType {
				return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument for toArray method must be a string"}
			}
			sep = args[0].(StringValue).Value
		}

		arr := ArrayValue{Value: make([]RuntimeValue, 0)}

		values := strings.Split(s.Value, sep)

		for _, value := range values {
			arr.Value = append(arr.Value, StringValue{Value: value})
		}

		return &arr
	}
}
func StringLen(s StringValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		return NumberValue{Value: float64(len(s.Value))}
	}
}
func StringSlice(s StringValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		length := len(s.Value)
		if len(args) == 2 {
			if args[0].GetType() != NumberType || args[1].GetType() != NumberType {
				return ErrorValue{ErrorType: InvalidArgumentError, Value: "Arguments for slice method should be numbers"}
			}
			from := int(args[0].(NumberValue).Value)
			to := int(args[1].(NumberValue).Value)
			if to < 0 {
				to = length + to
			}
			if from < 0 {
				from = length + from
			}
			if from >= length || to >= length || from > to {
				return ErrorValue{ErrorType: RuntimeError, Value: "Index out of range [" + strconv.Itoa(from) + ":" + strconv.Itoa(to) + "]"}
			}
			return StringValue{Value: s.Value[from:to]}
		} else if len(args) == 1 {
			if args[0].GetType() != NumberType {
				return ErrorValue{ErrorType: InvalidArgumentError, Value: "Argument for slice method should be a number"}
			}
			from := int(args[0].(NumberValue).Value)
			if from < 0 {
				from = length + from
			}
			if from >= length {
				return ErrorValue{ErrorType: RuntimeError, Value: "Index out of range [" + strconv.Itoa(from) + "]"}
			}
			return StringValue{Value: s.Value[from:]}
		} else {
			return ErrorValue{ErrorType: RuntimeError, Value: "Missing arguments for slice method, need at least one"}
		}
	}
}
