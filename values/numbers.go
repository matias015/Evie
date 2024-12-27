package values

import "fmt"

type NumberValue struct {
	Value float64
}

func (a NumberValue) GetString() string {
	return fmt.Sprintf("%v", a.Value)
}

func (a NumberValue) GetType() ValueType {
	return NumberType
}
func (b NumberValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
