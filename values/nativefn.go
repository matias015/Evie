package values

import "fmt"

type NativeFunctionValue struct {
	Value func(args []RuntimeValue) RuntimeValue
}

func (a NativeFunctionValue) GetNumber() float64 {
	return 1
}
func (a NativeFunctionValue) GetString() string {
	return "Native Function"
}
func (a NativeFunctionValue) GetBool() bool {
	return true
}
func (a NativeFunctionValue) GetType() ValueType {
	return NativeFunctionType
}
func (nfn NativeFunctionValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
