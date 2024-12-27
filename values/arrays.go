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

	return props[name], nil
}

func ArrayAdd(val *RuntimeValue, args []RuntimeValue) {
	// First, get the current slice from the interface{}
	currentSlice := (*val).(*ArrayValue)

	currentSlice.Value = append(currentSlice.Value, args...)

}
