package values

import "fmt"

type ErrorValue struct {
	Value string
}

func (a ErrorValue) GetBool() bool {
	return false
}
func (a ErrorValue) GetNumber() float64 {
	return 0
}
func (a ErrorValue) GetString() string {
	return fmt.Sprintf("%v", a.Value)
}
func (a ErrorValue) GetType() ValueType {
	return ErrorType
}
func (b ErrorValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
