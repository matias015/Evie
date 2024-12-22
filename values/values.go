package values

import (
	"evie/core"
	"evie/parser"
	"strconv"
)

type RuntimeValue interface {
	GetType() string
	GetStr() string
	GetNumber() float64
	GetBool() bool
	GetProp(*RuntimeValue, string) RuntimeValue
}

// This is kinda a mess
// GetProp allows to access to props of objects or native methods for native data types
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
	return true
}

func (s NumberValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	props := map[string]RuntimeValue{}

	props["isInteger"] = NativeFunctionValue{Value: func(args []RuntimeValue) RuntimeValue {
		return BooleanValue{Value: float64(int64(s.Value)) == s.Value}
	}}

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

func (s BooleanValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
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

func (s NativeFunctionValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
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
	Evaluator    core.Evaluator
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

func (s FunctionValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	return nil
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

func (a StructValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	return nil
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

func (a DictionaryValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	props := map[string]RuntimeValue{}

	props = map[string]RuntimeValue{

		"add": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].GetStr()
				value := args[1]
				a.Value[key] = value
				return BooleanValue{Value: true}
			},
		},
		"remove": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].GetStr()
				delete(a.Value, key)
				return BooleanValue{Value: true}
			},
		},
		"has": NativeFunctionValue{
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].GetStr()
				_, ok := a.Value[key]
				return BooleanValue{Value: ok}
			},
		},
	}

	return props[name]
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

func (a ObjectValue) GetProp(v *RuntimeValue, prop string) RuntimeValue {
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
	Exp   parser.Exp
	Env   interface{}
	Value RuntimeValue
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

func (s SignalValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	return SignalValue{Type: s.Type, Exp: s.Exp, Env: s.Env}
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

func (s ErrorValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	return ErrorValue{Value: s.Value}
}

type CapturedErrorValue struct {
	Value ErrorValue
}

func (s CapturedErrorValue) GetNumber() float64 {
	return 0
}
func (s CapturedErrorValue) GetType() string {
	return "CapturedErrorValue"
}

func (s CapturedErrorValue) GetStr() string {
	return s.Value.Value
}
func (s CapturedErrorValue) GetBool() bool {
	return false
}

func (s CapturedErrorValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	return ErrorValue{Value: s.Value.Value}
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

func (s NamespaceValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
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

func (s NothingValue) GetProp(v *RuntimeValue, name string) RuntimeValue {
	return nil
}
