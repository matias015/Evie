package values

type ObjectValue struct {
	Struct StructValue
	Value  map[string]RuntimeValue
}

func (a ObjectValue) GetString() string {
	return "Object"
}

func (a ObjectValue) GetType() ValueType {
	return ObjectType
}

func (a ObjectValue) GetProp(v *RuntimeValue, prop string) (RuntimeValue, error) {
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
