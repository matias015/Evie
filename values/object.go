package values

type ObjectValue struct {
	Struct StructValue
	Value  map[string]RuntimeValue
}

func (a ObjectValue) GetNumber() float64 {
	return 1
}
func (a ObjectValue) GetString() string {
	return "Object"
}
func (a ObjectValue) GetBool() bool {
	return true
}
func (a ObjectValue) GetType() ValueType {
	return ObjectType
}

func (a ObjectValue) GetProp(v *RuntimeValue, prop string) (RuntimeValue, error) {

	if prop == "get" {
		return NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				if len(args) < 1 {
					return ErrorValue{ErrorType: InvalidArgumentError, Value: "Missing arguments for get method"}
				}

				if args[0].GetType() != StringType {
					return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument for get method must be a string"}
				}

				key := args[0].(StringValue).Value
				if val, ok := a.Value[key]; ok {
					return val
				}

				return ErrorValue{ErrorType: RuntimeError, Value: "Key not found"}
			},
		}, nil
	}
	propValue, exists := a.Value[prop]
	if !exists {
		propValue = a.Struct.Methods[prop]
		switch propValue.GetType() {
		case NativeFunctionType:
			fn := propValue.(NativeFunctionValue)
			return fn, nil
		case FunctionType:
			fn := propValue.(FunctionValue)
			fn.StructObjRef = (*v).(*ObjectValue)
			return fn, nil
		}
	}

	return propValue, nil
}
