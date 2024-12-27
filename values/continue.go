package values

import "fmt"

type ContinueValue struct {
}

func (a ContinueValue) GetString() string {
	return "Continue"
}

// GetType returns the type of the value, which is BoolType for BoolValue
func (a ContinueValue) GetType() ValueType {
	return ContinueType
}
func (b ContinueValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
