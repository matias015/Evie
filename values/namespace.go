package values

import "fmt"

type NamespaceValue struct {
	Value map[string]RuntimeValue
}

func (a NamespaceValue) GetString() string {
	return "Namespace"
}

func (a NamespaceValue) GetType() ValueType {
	return NamespaceType
}
func (a NamespaceValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {

	prop, ex := a.Value[name]

	if !ex {
		return NothingValue{}, fmt.Errorf("property %s does not exists", name)
	}

	return prop, nil
}
