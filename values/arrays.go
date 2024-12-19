package values

import "strconv"

type ArrayValue struct {
	Value []RuntimeValue
}

func (a ArrayValue) GetNumber() float64 {
	return float64(len(a.Value))
}
func (a ArrayValue) GetType() string {
	return "ArrayValue"
}

func (a ArrayValue) GetStr() string {
	return "array"
}
func (a ArrayValue) GetBool() bool {
	if len(a.Value) > 0 {
		return true
	} else {
		return false
	}
}

func (a ArrayValue) GetProp(v *RuntimeValue, name string) RuntimeValue {

	props := map[string]RuntimeValue{
		"slice": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				length := len(a.Value)
				if len(args) == 2 {
					from := int(args[0].GetNumber())
					to := int(args[1].GetNumber())
					if to < 0 {
						to = length + to
					}
					if from < 0 {
						from = length + from
					}
					if from > length || to > length {
						return ErrorValue{Value: "Index out of range [" + strconv.FormatFloat(args[0].GetNumber(), 'f', -1, 64) + ":" + strconv.FormatFloat(args[1].GetNumber(), 'f', -1, 64) + "]"}
					}
					return &ArrayValue{Value: a.Value[from:to]}
				} else if len(args) == 1 {
					from := int(args[0].GetNumber())
					if from < 0 {
						from = length + from
					}
					if from > length {
						return ErrorValue{Value: "Index out of range [" + strconv.FormatFloat(args[0].GetNumber(), 'f', -1, 64) + "]"}
					}
					return &ArrayValue{Value: a.Value[from:]}
				} else {
					return &ArrayValue{Value: []RuntimeValue{}}
				}
			},
		},
		"add": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				val := (*v).(*ArrayValue)
				for _, arg := range args {
					val.Value = append(val.Value, arg)
				}
				return BooleanValue{Value: true}
			},
		},
		"addFirst": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				val := (*v).(*ArrayValue)
				for _, arg := range args {
					val.Value = append([]RuntimeValue{arg}, val.Value...)
				}
				return BooleanValue{Value: true}
			},
		},
		"has": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				if len(args) == 1 {
					for _, arg := range args {
						for _, val := range a.Value {
							if arg.GetStr() == val.GetStr() {
								if arg.GetType() == val.GetType() {
									return BooleanValue{Value: true}
								}
							}
						}
					}
					return BooleanValue{Value: false}
				} else if len(args) > 1 {

					has := false

					for _, arg := range args {
						has = false
						for _, val := range a.Value {
							if arg.GetStr() == val.GetStr() {
								if arg.GetType() == val.GetType() {
									has = true
									break
								}
							}
						}

						if has == false {
							return BooleanValue{Value: false}
						}

					}
					return BooleanValue{Value: true}
				} else {
					return BooleanValue{Value: false}
				}
			},
		},
		"find": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				arg := args[0]

				for index, value := range a.Value {
					if arg.GetStr() == value.GetStr() {
						if arg.GetType() == value.GetType() {
							return NumberValue{Value: float64(index)}
						}
					}
				}

				return NumberValue{Value: float64(-1)}
			},
		},
		"len": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				return NumberValue{Value: float64(len(a.Value))}
			},
		},
	}

	return props[name]
}
