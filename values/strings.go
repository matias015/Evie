package values

/*
----------------------------------------------------------
--- StringValue
----------------------------------------------------------
*/
// type StringValue struct {
// 	Value   string
// 	Mutable bool
// }

// func (s StringValue) GetStr() string     { return s.Value }
// func (s StringValue) GetNumber() float64 { return 0 }
// func (s StringValue) GetBool() bool {
// 	if s.Value == "" {
// 		return false
// 	} else {
// 		return true
// 	}
// }
// func (s StringValue) GetType() string {
// 	return "StringValue"
// }

// func (s StringValue) GetProp(v *RuntimeValue, name string) RuntimeValue {

// 	props := map[string]RuntimeValue{

//
//
// 		"is": NativeFunctionValue{
// 			Value: func(args []RuntimeValue) RuntimeValue {
// 				for _, arg := range args {
// 					if arg.GetStr() == s.GetStr() {
// 						return BooleanValue{Value: true}
// 					}
// 				}
// 				return BooleanValue{Value: false}
// 			},
// 		},
// 		"addPaddingLeft": NativeFunctionValue{
// 			Value: func(args []RuntimeValue) RuntimeValue {
// 				char := args[0].GetStr()
// 				length := int(args[1].GetNumber())
// 				actualLength := len(s.GetStr())
// 				output := s.Value
// 				for i := 0; i < length-actualLength; i++ {
// 					output = char + output
// 				}
// 				return StringValue{Value: output}
// 			},
// 		},
// 		"addPaddingRight": NativeFunctionValue{
// 			Value: func(args []RuntimeValue) RuntimeValue {
// 				char := args[0].GetStr()
// 				length := int(args[1].GetNumber())
// 				actualLength := len(s.GetStr())
// 				output := s.Value
// 				for i := 0; i < length-actualLength; i++ {
// 					output = output + char
// 				}
// 				return StringValue{Value: output}
// 			},
// 		},
// 		"trim": NativeFunctionValue{
// 			Value: func(args []RuntimeValue) RuntimeValue {

// 				needed := " "
// 				if len(args) > 0 {
// 					needed = args[0].GetStr()
// 				}
// 				return StringValue{Value: strings.Trim(s.Value, needed)}
// 			},
// 		},
// 		"toArray": NativeFunctionValue{
// 			Value: func(args []RuntimeValue) RuntimeValue {

// 				sep := ""

// 				if len(args) > 0 {
// 					sep = args[0].GetStr()
// 				}

// 				arr := ArrayValue{Value: make([]RuntimeValue, 0)}

// 				values := strings.Split(s.Value, sep)

// 				for _, value := range values {
// 					arr.Value = append(arr.Value, StringValue{Value: value})
// 				}

// 				return &arr
// 			},
// 		},
// 		"slice": NativeFunctionValue{
// 			Value: func(args []RuntimeValue) RuntimeValue {
// 				length := len(s.Value)
// 				if len(args) == 2 {
// 					from := int(args[0].GetNumber())
// 					to := int(args[1].GetNumber())
// 					if to < 0 {
// 						to = length + to
// 					}
// 					if from < 0 {
// 						from = length + from
// 					}
// 					if from > length || to > length {
// 						return ErrorValue{Value: "Index out of range [" + strconv.FormatFloat(args[0].GetNumber(), 'f', -1, 64) + ":" + strconv.FormatFloat(args[1].GetNumber(), 'f', -1, 64) + "]"}
// 					}
// 					return StringValue{Value: s.Value[from:to]}
// 				} else if len(args) == 1 {
// 					from := int(args[0].GetNumber())
// 					if from < 0 {
// 						from = length + from
// 					}
// 					if from > length {
// 						return ErrorValue{Value: "Index out of range [" + strconv.FormatFloat(args[0].GetNumber(), 'f', -1, 64) + "]"}
// 					}
// 					return StringValue{Value: s.Value[from:]}
// 				} else {
// 					return StringValue{Value: ""}
// 				}
// 			},
// 		},
// 	}

// 	return props[name]
// }

func StringLength(v *RuntimeValue) RuntimeValue {
	return RuntimeValue{Type: NumberType, Value: float64(len(v.Value.(string)))}
}
