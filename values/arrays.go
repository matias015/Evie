package values

import (
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

func (a *ArrayValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {

	props := map[string]RuntimeValue{
		"slice": NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
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
					return ErrorValue{Value: "Index out of range [" + strconv.Itoa(from) + ":" + strconv.Itoa(to) + "]"}
				}
				return &ArrayValue{Value: a.Value[from:to]}
			} else if len(args) == 1 {
				from := int(args[0].(NumberValue).Value)
				if from < 0 {
					from = length + from
				}
				if from > length {
					return ErrorValue{Value: "Index out of range [" + strconv.Itoa(from) + "]"}
				}
				return &ArrayValue{Value: a.Value[from:]}
			} else {
				return &ArrayValue{Value: []RuntimeValue{}}
			}
		},
		},
		"add": NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
			for _, arg := range args {
				a.Value = append(a.Value, arg)
			}
			return BoolValue{Value: true}
		},
		},
		"addFirst": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				for _, arg := range args {
					a.Value = append([]RuntimeValue{arg}, a.Value...)
				}
				return BoolValue{Value: true}
			},
		},
		"has": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
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
			},
		},
		"find": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {

				if len(args) == 0 {
					return ErrorValue{Value: "Missing argument for find method"}
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
			},
		},
		"len": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				return NumberValue{Value: float64(len(a.Value))}
			},
		},
	}

	props["remove"] = NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
		if len(args) == 1 {
			if args[0].GetType() != NumberType {
				return ErrorValue{Value: "First argument must be a number"}
			}

			index := int(args[0].(NumberValue).Value)

			if index < 0 {
				index = len(a.Value) + index
			}

			if index >= 0 && index < len(a.Value) {
				a.Value = append(a.Value[:index], a.Value[index+1:]...)
				return BoolValue{Value: true}
			}

			return ErrorValue{Value: "Index out of range or array is empty"}
		} else {
			a.Value = a.Value[:len(a.Value)-1]
		}
		return NothingValue{}
	}}

	props["removeFirst"] = NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
		if len(a.Value) == 0 {
			return ErrorValue{Value: "Array is empty"}
		}
		a.Value = a.Value[1:]
		return NothingValue{}
	}}

	props["addPaddingLeft"] = NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
		if len(args) < 2 {
			return ErrorValue{Value: "Missing arguments for addPaddingLeft method"}
		}

		if args[1].GetType() != NumberType {
			return ErrorValue{Value: "Second argument for addPaddingLeft method must be a number"}
		}

		char := args[0].GetString()
		length := int(args[1].(NumberValue).Value)
		output := make([]RuntimeValue, length)
		for i := range output {
			output[i] = StringValue{Value: char}
		}
		output = append(output, a.Value...)
		return &ArrayValue{Value: output}
	}}

	props["addPaddingRight"] = NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
		if len(args) < 2 {
			return ErrorValue{Value: "Missing arguments for addPaddingRight method"}
		}

		if args[1].GetType() != NumberType {
			return ErrorValue{Value: "Second argument for addPaddingRight method must be a number"}
		}

		char := args[0].GetString()
		length := int(args[1].(NumberValue).Value)
		output := make([]RuntimeValue, length)
		for i := range output {
			output[i] = StringValue{Value: char}
		}
		output = append(a.Value, output...)
		return &ArrayValue{Value: output}
	}}

	return props[name], nil
}

func ArrayAdd(val *RuntimeValue, args []RuntimeValue) {
	// First, get the current slice from the interface{}
	currentSlice := (*val).(*ArrayValue)

	currentSlice.Value = append(currentSlice.Value, args...)

}
