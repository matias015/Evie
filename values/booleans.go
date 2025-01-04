package values

import "fmt"

type BoolValue struct {
	Value bool
}

func (a BoolValue) GetString() string {
	return fmt.Sprintf("%v", a.Value)
}

func (a BoolValue) GetBool() bool {
	return a.Value
}
func (a BoolValue) GetNumber() float64 {
	if a.Value {
		return 1
	}
	return 0
}

// GetType returns the type of the value, which is BoolType for BoolValue
func (a BoolValue) GetType() ValueType {
	return BoolType
}
func (b BoolValue) GetProp(name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
