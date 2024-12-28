package values

import "fmt"

type ContinueValue struct {
}

func (a ContinueValue) GetString() string {
	return "Continue"
}

func (a ContinueValue) GetBool() bool {
	return false
}
func (a ContinueValue) GetNumber() float64 {
	return 0
}

// GetType returns the type of the value, which is BoolType for BoolValue
func (a ContinueValue) GetType() ValueType {
	return ContinueType
}
func (b ContinueValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
