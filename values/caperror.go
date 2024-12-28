package values

import "fmt"

type CapturedErrorValue struct {
	Value string
}

func (a CapturedErrorValue) GetBool() bool {
	return false
}
func (a CapturedErrorValue) GetNumber() float64 {
	return 0
}
func (a CapturedErrorValue) GetString() string {
	return fmt.Sprintf("%v", a.Value)
}
func (a CapturedErrorValue) GetType() ValueType {
	return CapturedErrorType
}
func (b CapturedErrorValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
