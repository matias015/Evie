package values

import (
	"evie/core"
	"evie/parser"
	"fmt"
	"strconv"
)

type ValueType uint8

const (
	StringType ValueType = iota
	NumberType
	BoolType
	NothingType
	ErrorType
	StructType
	DictionaryType
	ContinueType
	BreakType
	FunctionType
	NativeFunctionType
	ArrayType
	ReturnType
	ObjectType
	NamespaceType
	CapturedErrorType
)

func (v ValueType) String() string {
	return [...]string{
		"string",
		"number",
		"boolean",
		"nothing",
		"error",
		"StructType",
		"DictionaryType",
		"ContinueType",
		"BreakType",
		"FunctionType",
		"NativeFunctionType",
		"array",
		"ReturnType",
		"ObjectType",
		"NamespaceType",
		"CapturedErrorType",
	}[v]
}

type RuntimeValue struct {
	Type  ValueType
	Value interface{}
}

func (val *RuntimeValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	props := map[string]map[string]RuntimeValue{}

	if val.Type == NamespaceType {
		prop, exists := val.Value.(map[string]RuntimeValue)[name]

		if !exists {
			return RuntimeValue{}, fmt.Errorf("property %s does not exists", name)
		}

		return prop, nil
	}

	if val.Type == ObjectType {
		return val.Value.(*ObjectValue).GetProp(v, name)
	}

	if val.Type == DictionaryType {
		return val.Value.(*DictionaryValue).GetProp(v, name)
	}

	props["string"] = map[string]RuntimeValue{
		"len": {Type: NativeFunctionType, Value: func(args []RuntimeValue) RuntimeValue {
			return StringLength(v)
		}},
		"slice": {Type: NativeFunctionType, Value: func(args []RuntimeValue) RuntimeValue {
			value := val.Value.(string)
			length := len(value)
			if len(args) == 2 {
				from := int(args[0].Value.(float64))
				to := int(args[1].Value.(float64))
				if to < 0 {
					to = length + to
				}
				if from < 0 {
					from = length + from
				}
				if from > length || to > length {
					return RuntimeValue{Type: ErrorType, Value: "Index out of range [" + strconv.FormatFloat(args[0].Value.(float64), 'f', -1, 64) + ":" + strconv.FormatFloat(args[1].Value.(float64), 'f', -1, 64) + "]"}
				}
				return RuntimeValue{Type: StringType, Value: value[from:to]}
			} else if len(args) == 1 {
				from := int(args[0].Value.(float64))
				if from < 0 {
					from = length + from
				}
				if from > length {
					return RuntimeValue{Type: ErrorType, Value: "Index out of range [" + strconv.FormatFloat(args[0].Value.(float64), 'f', -1, 64) + "]"}
				}
				return RuntimeValue{Type: StringType, Value: value[from:]}
			} else {
				return RuntimeValue{Type: StringType, Value: ""}
			}
		}},
	}

	if val.Type == ArrayType {
		value := val.Value.(*ArrayValue)
		fn, err := value.GetProp(v.Value.(*ArrayValue), name)
		if err != nil {
			return RuntimeValue{}, err
		}
		return fn, nil
	}

	prop, exists := props[val.Type.String()][name]

	if !exists {
		return RuntimeValue{}, fmt.Errorf("property %s does not exists", name)
	}

	return prop, nil
}

type NativeFunctionValue struct {
	Value func(args []RuntimeValue) RuntimeValue
}

/*
----------------------------------------------------------
--- FUCNTION
----------------------------------------------------------
*/

type FunctionValue struct {
	Struct       string
	StructObjRef *ObjectValue
	Body         []parser.Stmt
	Parameters   []string
	Environment  interface{}
	Evaluator    core.Evaluator
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

/*
----------------------------------------------------------
--- DICTIONARIES
----------------------------------------------------------
*/

func (a *DictionaryValue) GetProp(v *RuntimeValue, name string) (RuntimeValue, error) {
	props := map[string]RuntimeValue{

		"add": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].Value.(string)
				value := args[1]
				a.Value[key] = value
				return RuntimeValue{Type: BoolType, Value: true}
			},
		},
		"remove": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].Value.(string)
				delete(a.Value, key)
				return RuntimeValue{Type: BoolType, Value: true}
			},
		},
		"has": {Type: NativeFunctionType,
			Value: func(args []RuntimeValue) RuntimeValue {
				key := args[0].Value.(string)
				_, ok := a.Value[key]
				return RuntimeValue{Type: BoolType, Value: ok}
			},
		},
	}

	prop, ex := props[name]

	if !ex {
		return RuntimeValue{}, fmt.Errorf("property %s does not exists", name)
	}

	return prop, nil
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

func (a ObjectValue) GetProp(v *RuntimeValue, prop string) (RuntimeValue, error) {
	propValue, exists := a.Value[prop]

	if !exists {
		propValue = a.Struct.Methods[prop]
		switch propValue.Type {
		case NativeFunctionType:
			fn := propValue.Value.(NativeFunctionValue)
			return RuntimeValue{Type: NativeFunctionType, Value: fn}, nil
		case FunctionType:
			fn := propValue.Value.(FunctionValue)
			fn.StructObjRef = v.Value.(*ObjectValue)
			return RuntimeValue{Type: FunctionType, Value: fn}, nil
		}
	}

	return propValue, nil
}

type DictionaryValue struct {
	Value map[string]RuntimeValue
}
