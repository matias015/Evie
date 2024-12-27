package values

import "fmt"

type BreakValue struct{}

func (a BreakValue) GetString() string {
	return "BreakValue"
}

// GetType returns the type of the value, which is BoolType for BoolValue
func (a BreakValue) GetType() ValueType {
	return BreakType
}
func (b BreakValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
