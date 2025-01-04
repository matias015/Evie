package values

import (
	"fmt"
	"math"
)

type NumberValue struct {
	Value float64
}

func (a NumberValue) GetNumber() float64 {
	return a.Value
}
func (a NumberValue) GetString() string {
	return fmt.Sprintf("%v", a.Value)
}
func (a NumberValue) GetBool() bool {
	return a.Value != 0
}
func (a NumberValue) GetType() ValueType {
	return NumberType
}
func (n NumberValue) GetProp(name string) (RuntimeValue, error) {

	if name == "round" {
		return NativeFunctionValue{Value: NumberRound(n)}, nil
	} else if name == "isInt" {
		return NativeFunctionValue{Value: NumberIsInt(n)}, nil
	} else {
		return nil, fmt.Errorf("Property %s not found", name)
	}

}

func NumberRound(n NumberValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		mult := 1.0

		if len(args) >= 1 {
			if args[0].GetType() != NumberType {
				return ErrorValue{ErrorType: InvalidArgumentError, Value: "First argument must be a number"}
			}
			mult = args[0].(NumberValue).Value
			if mult < 1 {
				mult = 1
			}
		}

		val := math.Round(n.Value*mult) / mult
		return NumberValue{Value: val}
	}
}

func NumberIsInt(n NumberValue) func([]RuntimeValue) RuntimeValue {
	return func(args []RuntimeValue) RuntimeValue {
		return BoolValue{Value: n.Value == float64(int(n.Value))}
	}
}
