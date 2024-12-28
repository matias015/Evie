package values

import "fmt"

type NothingValue struct {
}

func (a NothingValue) GetNumber() float64 {
	return 0
}
func (a NothingValue) GetString() string {
	return "Nothing"
}
func (a NothingValue) GetBool() bool {
	return false
}
func (a NothingValue) GetType() ValueType {
	return NothingType
}
func (b NothingValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
