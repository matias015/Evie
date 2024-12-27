package values

import "fmt"

type StructValue struct {
	Properties []string
	Methods    map[string]RuntimeValue
}

func (a StructValue) GetString() string {
	return "Struct"
}

func (a StructValue) GetType() ValueType {
	return StructType
}
func (s StructValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
