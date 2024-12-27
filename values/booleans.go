package values

import "fmt"

type BoolValue struct {
	Value bool
}

func (a BoolValue) GetString() string {
	return fmt.Sprintf("%v", a.Value)
}

// GetType returns the type of the value, which is BoolType for BoolValue
func (a BoolValue) GetType() ValueType {
	return BoolType
}
func (b BoolValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
