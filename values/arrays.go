package values

import (
	"fmt"
	"strconv"
)

type ArrayValue struct {
	Value []RuntimeValue
}

func (a ArrayValue) GetType() ValueType {
	return ArrayType
}

func (a ArrayValue) GetString() string {
	return "array"
}
func (a ArrayValue) GetNumber() float64 {
	return 0
}
func (a ArrayValue) GetBool() bool {
	return len(a.Value) > 0
}

func (a *ArrayValue) GetProp(name string) (RuntimeValue, error) {

	if name == "add" {
		return NativeFunctionValue{Value: ArrayAdd(a)}, nil
	} else if name == "len" {
		return NativeFunctionValue{Value: ArrayLength(a)}, nil
	} else if name == "slice" {
		return NativeFunctionValue{Value: ArraySlice(a)}, nil
	} else if name == "has" {
		return NativeFunctionValue{Value: ArrayHas(a)}, nil
	} else if name == "find" {
		return NativeFunctionValue{Value: ArrayFind(a)}, nil
	} else if name == "remove" {
		return NativeFunctionValue{Value: ArrayRemove(a)}, nil
	} else if name == "addPaddingLeft" {
		return NativeFunctionValue{Value: ArrayAddPaddingLeft(a)}, nil
	} else if name == "addPaddingRight" {
		return NativeFunctionValue{Value: ArrayAddPaddingRight(a)}, nil
	} else if name == "addFirst" {
		return NativeFunctionValue{Value: ArrayAddFirst(a)}, nil
	} else if name == "removeFirst" {
		return NativeFunctionValue{Value: ArrayRemoveFirst(a)}, nil
	} else {
		return NothingValue{}, fmt.Errorf("property %s does not exists", name)
	}
}

func ArrayAddFirst(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		for _, arg := range args {
			a.Value = append([]RuntimeValue{arg}, a.Value...)
		}
		return BoolValue{Value: true}
	}
}

func ArrayRemoveFirst(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(a.Value) == 0 {
			return ErrorValue{ErrorType: RuntimeError, Value: "Array is empty"}
		}
		a.Value = a.Value[1:]
		return NothingValue{}
	}
}

func ArrayAdd(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		for _, arg := range args {
			a.Value = append(a.Value, arg)
		}
		return BoolValue{Value: true}
	}
}
func ArrayLength(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		return NumberValue{Value: float64(len(a.Value))}
	}
}

func ArraySlice(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		length := len(a.Value)
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
				return ErrorValue{ErrorType: InvalidIndexError, Value: "Index out of range [" + strconv.Itoa(from) + ":" + strconv.Itoa(to) + "]"}
			}
			return &ArrayValue{Value: a.Value[from:to]}
		} else if len(args) == 1 {
			from := int(args[0].(NumberValue).Value)
			if from < 0 {
				from = length + from
			}
			if from > length {
				return ErrorValue{ErrorType: RuntimeError, Value: "Index out of range [" + strconv.Itoa(from) + "]"}
			}
			return &ArrayValue{Value: a.Value[from:]}
		} else {
			return &ArrayValue{Value: []RuntimeValue{}}
		}

	}
}

func ArrayRemove(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) == 1 {
			if args[0].GetType() != NumberType {
				return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument must be a number"}
			}

			index := int(args[0].(NumberValue).Value)

			if index < 0 {
				index = len(a.Value) + index
			}

			if index >= 0 && index < len(a.Value) {
				a.Value = append(a.Value[:index], a.Value[index+1:]...)
				return BoolValue{Value: true}
			}

			return ErrorValue{ErrorType: InvalidIndexError, Value: "Index out of range or array is empty"}
		} else {
			a.Value = a.Value[:len(a.Value)-1]
		}
		return NothingValue{}
	}
}

func ArrayHas(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {

		if len(args) == 1 {
			for _, arg := range args {
				for _, val := range a.Value {
					if arg.GetType() == val.GetType() {
						if arg.GetString() == val.GetString() {
							return BoolValue{Value: true}
						}
					}
				}
			}
			return BoolValue{Value: false}
		} else if len(args) > 1 {

			has := false

			for _, arg := range args {
				has = false
				for _, val := range a.Value {
					if arg.GetType() == val.GetType() {
						if arg.GetString() == val.GetString() {
							has = true
							break
						}
					}
				}

				if has == false {
					return BoolValue{Value: false}
				}

			}
			return BoolValue{Value: true}
		} else {
			return BoolValue{Value: false}
		}
	}
}

func ArrayFind(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) == 0 {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing argument for find method"}
		}

		arg := args[0]

		for index, value := range a.Value {
			if arg.GetType() == value.GetType() {
				if arg.GetString() == value.GetString() {
					return NumberValue{Value: float64(index)}
				}
			}
		}

		return NumberValue{Value: -1.0}
	}
}

func ArrayAddPaddingLeft(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) < 2 {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing arguments for addPaddingLeft method"}
		}

		if args[1].GetType() != NumberType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Second argument for addPaddingLeft method must be a number"}
		}

		char := args[0].GetString()
		length := int(args[1].(NumberValue).Value)
		output := make([]RuntimeValue, length)
		for i := range output {
			output[i] = StringValue{Value: char}
		}
		output = append(output, a.Value...)
		return &ArrayValue{Value: output}
	}
}

func ArrayAddPaddingRight(a *ArrayValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		if len(args) < 2 {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing arguments for addPaddingRight method"}
		}

		if args[1].GetType() != NumberType {
			return ErrorValue{ErrorType: InvalidArgumentError, Value: "Second argument for addPaddingRight method must be a number"}
		}

		char := args[0].GetString()
		length := int(args[1].(NumberValue).Value)
		output := make([]RuntimeValue, length)
		for i := range output {
			output[i] = StringValue{Value: char}
		}
		output = append(a.Value, output...)
		return &ArrayValue{Value: output}
	}
}
