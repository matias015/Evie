package values

import (
	"evie/parser"
	"strconv"
)

type RuntimeValue interface {
	GetType() string
	GetStr() string
	GetNumber() float64
	GetBool() bool
	GetProp(string) RuntimeValue
}

/*
----------------------------------------------------------
--- Number value
----------------------------------------------------------
*/
// NUMBER VALUE
type NumberValue struct {
	Value float64
}

func (s NumberValue) GetNumber() float64 {
	return s.Value
}
func (s NumberValue) GetType() string {
	return "NumberValue"
}

func (s NumberValue) GetStr() string {
	return strconv.FormatFloat(s.Value, 'f', -1, 64)
}
func (s NumberValue) GetBool() bool {
	return false
}

func (s NumberValue) GetProp(name string) RuntimeValue {
	props := map[string]RuntimeValue{}

	props = map[string]RuntimeValue{
		"toString": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				return StringValue{Value: s.GetStr()}
			},
		},
	}

	return props[name]
}

/*
----------------------------------------------------------
--- Boolean value
----------------------------------------------------------
*/
type BooleanValue struct {
	Value bool
}

func (s BooleanValue) GetNumber() float64 {
	if s.Value {
		return 1
	} else {
		return 0
	}
}
func (s BooleanValue) GetType() string {
	return "BooleanValue"
}

func (s BooleanValue) GetStr() string {
	if s.Value {
		return "true"
	} else {
		return "false"
	}
}

func (s BooleanValue) GetBool() bool {
	return s.Value
}

func (s BooleanValue) GetProp(name string) RuntimeValue {
	return BooleanValue{Value: s.Value}
}

/*
----------------------------------------------------------
--- Boolean value
----------------------------------------------------------
*/

type NativeFunctionValue struct {
	StructObjRef ObjectValue
	Value        func(args []RuntimeValue) RuntimeValue
}

func (s NativeFunctionValue) GetNumber() float64 {
	return 1
}
func (s NativeFunctionValue) GetType() string {
	return "NativeFunctionValue"
}

func (s NativeFunctionValue) GetStr() string {
	return "NativeFunctionValue"
}
func (s NativeFunctionValue) GetBool() bool {
	return true
}

func (s NativeFunctionValue) GetProp(name string) RuntimeValue {
	return NativeFunctionValue{Value: s.Value}
}

/*
----------------------------------------------------------
--- FUCNTION
----------------------------------------------------------
*/

type FunctionValue struct {
	Struct       string
	StructObjRef ObjectValue
	Body         []parser.Stmt
	Parameters   []string
	Environment  interface{}
}

func (s FunctionValue) GetNumber() float64 {
	return 1
}
func (s FunctionValue) GetType() string {
	return "FunctionValue"
}

func (s FunctionValue) GetStr() string {
	return "FunctionValue"
}
func (s FunctionValue) GetBool() bool {
	return true
}

func (s FunctionValue) GetProp(name string) RuntimeValue {
	return FunctionValue{Body: s.Body, Parameters: s.Parameters}
}

/*
----------------------------------------------------------
--- FUCNTION
----------------------------------------------------------
*/

type ArrayValue struct {
	Value []RuntimeValue
}

func (a ArrayValue) GetNumber() float64 {
	return float64(len(a.Value))
}
func (a ArrayValue) GetType() string {
	return "ArrayValue"
}

func (a ArrayValue) GetStr() string {
	return "array"
}
func (a ArrayValue) GetBool() bool {
	if len(a.Value) > 0 {
		return true
	} else {
		return false
	}
}

func (a ArrayValue) GetProp(name string) RuntimeValue {
	props := map[string]RuntimeValue{}

	props = map[string]RuntimeValue{
		"slice": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				length := len(a.Value)
				if len(args) == 2 {
					from := int(args[0].GetNumber())
					to := int(args[1].GetNumber())
					if to < 0 {
						to = length + to
					}
					if from < 0 {
						from = length + from
					}
					if from > length || to > length {
						return ErrorValue{Value: "Index out of range [" + strconv.FormatFloat(args[0].GetNumber(), 'f', -1, 64) + ":" + strconv.FormatFloat(args[1].GetNumber(), 'f', -1, 64) + "]"}
					}
					return ArrayValue{Value: a.Value[from:to]}
				} else if len(args) == 1 {
					from := int(args[0].GetNumber())
					if from < 0 {
						from = length + from
					}
					if from > length {
						return ErrorValue{Value: "Index out of range [" + strconv.FormatFloat(args[0].GetNumber(), 'f', -1, 64) + "]"}
					}
					return ArrayValue{Value: a.Value[from:]}
				} else {
					return ArrayValue{Value: []RuntimeValue{}}
				}
			},
		},
	}

	return props[name]
}

/*
----------------------------------------------------------
--- STRUCT
----------------------------------------------------------
*/

type StructValue struct {
	Properties []string
	Methods    map[string]RuntimeValue
}

func (a StructValue) GetNumber() float64 {
	return 1
}
func (a StructValue) GetType() string {
	return "StructValue"
}

func (a StructValue) GetStr() string {
	return "StructValue"
}
func (a StructValue) GetBool() bool {
	return true
}

func (a StructValue) GetProp(name string) RuntimeValue {
	return StructValue{Properties: a.Properties}
}

/*
----------------------------------------------------------
--- DICTIONARIES
----------------------------------------------------------
*/

type DictionaryValue struct {
	Value map[string]RuntimeValue
}

func (a DictionaryValue) GetNumber() float64 {
	return 1
}
func (a DictionaryValue) GetType() string {
	return "DictionaryValue"
}

func (a DictionaryValue) GetStr() string {
	return "DictionaryValue"
}
func (a DictionaryValue) GetBool() bool {
	return true
}

func (a DictionaryValue) GetProp(name string) RuntimeValue {
	return DictionaryValue{Value: a.Value}
}

/*
----------------------------------------------------------
--- OBJECT
----------------------------------------------------------
*/

type ObjectValue struct {
	Struct StructValue
	Value  map[string]RuntimeValue
}

func (a ObjectValue) GetNumber() float64 {
	return 1
}
func (a ObjectValue) GetType() string {
	return "ObjectValue"
}

func (a ObjectValue) GetStr() string {
	return "ObjectValue"
}
func (a ObjectValue) GetBool() bool {
	return true
}

func (a ObjectValue) GetProp(prop string) RuntimeValue {
	propValue := a.Value[prop]

	if propValue == nil {
		propValue = a.Struct.Methods[prop]
		switch fn := propValue.(type) {
		case NativeFunctionValue:
			fn.StructObjRef = a
			return fn
		case FunctionValue:
			fn.StructObjRef = a
			return fn
		}
	}

	return propValue
}

// Represents values that are used to break/continue/return
// Save the environment to calculate return value in functions that
// was declared inside other blocks like if, for, etc.
type SignalValue struct {
	Type  string
	Value parser.Exp
	Env   interface{}
}

func (s SignalValue) GetNumber() float64 {
	return 1
}
func (s SignalValue) GetType() string {
	return s.Type
}

func (s SignalValue) GetStr() string {
	return "SignalValue"
}
func (s SignalValue) GetBool() bool {
	return true
}

func (s SignalValue) GetProp(name string) RuntimeValue {
	return SignalValue{Type: s.Type, Value: s.Value, Env: s.Env}
}

type ErrorValue struct {
	Value string
}

func (s ErrorValue) GetNumber() float64 {
	return 1
}
func (s ErrorValue) GetType() string {
	return "ErrorValue"
}

func (s ErrorValue) GetStr() string {
	return s.Value
}
func (s ErrorValue) GetBool() bool {
	return true
}

func (s ErrorValue) GetProp(name string) RuntimeValue {
	return ErrorValue{Value: s.Value}
}

type NamespaceValue struct {
	Value map[string]RuntimeValue
}

func (s NamespaceValue) GetNumber() float64 {
	return 1
}
func (s NamespaceValue) GetType() string {
	return "NamespaceValue"
}

func (s NamespaceValue) GetStr() string {
	return "NamespaceValue"
}
func (s NamespaceValue) GetBool() bool {
	return true
}

func (s NamespaceValue) GetProp(name string) RuntimeValue {
	return s.Value[name]
}

type NothingValue struct {
}

func (s NothingValue) GetNumber() float64 {
	return 0
}
func (s NothingValue) GetType() string {
	return "NothingValue"
}

func (s NothingValue) GetStr() string {
	return "Nothing"
}

func (s NothingValue) GetBool() bool {
	return false
}

func (s NothingValue) GetProp(name string) RuntimeValue {
	return NothingValue{}
}
