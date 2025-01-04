package values

import "fmt"

const (
	RuntimeError           string = "RuntimeError"
	TypeError              string = "TypeError"
	InvalidIndexError      string = "InvalidIndexError"
	IdentifierError        string = "IdentifierError"
	ZeroDivisionError      string = "ZeroDivisionError"
	InvalidArgumentError   string = "InvalidArgumentError"
	InvalidConversionError string = "InvalidConversionError"
	CircularImportError    string = "CircularImportError"
	PropertyError          string = "PropertyError"
)

type ErrorValue struct {
	Value     string
	ErrorType string
	Object    *ObjectValue
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
func (b ErrorValue) GetProp(name string) (RuntimeValue, error) {
	methods := make(map[string]RuntimeValue, 2)

	// methods

	prop, ex := methods[name]

	if !ex {
		return NothingValue{}, fmt.Errorf("property %s does not exists", name)
	}

	return prop, nil

}
