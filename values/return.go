package values

import "fmt"

type ReturnValue struct {
	Value RuntimeValue
}

func (a ReturnValue) GetNumber() float64 {
	return 1
}
func (a ReturnValue) GetString() string {
	return "Return"
}
func (a ReturnValue) GetBool() bool {
	return true
}

// GetType returns the type of the value, which is BoolType for BoolValue
func (a ReturnValue) GetType() ValueType {
	return ReturnType
}
func (b ReturnValue) GetProp(name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
