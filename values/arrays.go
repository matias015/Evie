package values

import "strconv"

type ArrayValue struct {
	Value []RuntimeValue
}

func (a *ArrayValue) GetProp(v *ArrayValue, name string) (RuntimeValue, error) {

	props := map[string]RuntimeValue{
		"slice": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				length := len(a.Value)
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
					return RuntimeValue{Type: ArrayType, Value: &ArrayValue{Value: a.Value[from:to]}}
				} else if len(args) == 1 {
					from := int(args[0].Value.(float64))
					if from < 0 {
						from = length + from
					}
					if from > length {
						return RuntimeValue{Type: ErrorType, Value: "Index out of range [" + strconv.FormatFloat(args[0].Value.(float64), 'f', -1, 64) + "]"}
					}
					return RuntimeValue{Type: ArrayType, Value: &ArrayValue{Value: a.Value[from:]}}
				} else {
					return RuntimeValue{Type: ArrayType, Value: &ArrayValue{Value: []RuntimeValue{}}}
				}
			},
		},
		"add": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				for _, arg := range args {
					a.Value = append(a.Value, arg)
				}
				return RuntimeValue{Type: BoolType, Value: true}
			},
		},
		"addFirst": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				for _, arg := range args {
					a.Value = append([]RuntimeValue{arg}, a.Value...)
				}
				return RuntimeValue{Type: BoolType, Value: true}
			},
		},
		// "has": RuntimeValue{Type: NativeFunctionType,
		// 	Value: func(args []RuntimeValue) RuntimeValue {
		// 		if len(args) == 1 {
		// 			for _, arg := range args {
		// 				for _, val := range a.Value {
		// 					if arg.GetStr() == val.GetStr() {
		// 						if arg.GetType() == val.GetType() {
		// 							return BooleanValue{Value: true}
		// 						}
		// 					}
		// 				}
		// 			}
		// 			return BooleanValue{Value: false}
		// 		} else if len(args) > 1 {

		// 			has := false

		// 			for _, arg := range args {
		// 				has = false
		// 				for _, val := range a.Value {
		// 					if arg.GetStr() == val.GetStr() {
		// 						if arg.GetType() == val.GetType() {
		// 							has = true
		// 							break
		// 						}
		// 					}
		// 				}

		// 				if has == false {
		// 					return BooleanValue{Value: false}
		// 				}

		// 			}
		// 			return BooleanValue{Value: true}
		// 		} else {
		// 			return BooleanValue{Value: false}
		// 		}
		// 	},
		// },
		// "find": RuntimeValue{Type: NativeFunctionType,
		// 	Value: func(args []RuntimeValue) RuntimeValue {
		// 		arg := args[0]

		// 		for index, value := range a.Value {
		// 			if arg.GetStr() == value.GetStr() {
		// 				if arg.GetType() == value.GetType() {
		// 					return NumberValue{Value: float64(index)}
		// 				}
		// 			}
		// 		}

		// 		return NumberValue{Value: float64(-1)}
		// 	},
		// },
		"len": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				return RuntimeValue{Type: NumberType, Value: float64(len(a.Value))}
			},
		},
	}

	return props[name], nil
}

func ArrayAdd(val *RuntimeValue, args []RuntimeValue) {
	// First, get the current slice from the interface{}
	currentSlice := val.Value.(*[]RuntimeValue)

	// Append the new value
	newSlice := append(*currentSlice, RuntimeValue{Type: NumberType, Value: float64(34)})

	// Important: Update the interface{} Value with the new slice
	val.Value = newSlice
}
