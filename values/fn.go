package values

import (
	"evie/core"
	"evie/parser"
	"fmt"
)

type FunctionValue struct {
	Struct       string
	StructObjRef *ObjectValue
	Body         []parser.Stmt
	Parameters   []string
	Environment  interface{}
	Evaluator    core.Evaluator
}

func (a FunctionValue) GetNumber() float64 {
	return 1
}
func (a FunctionValue) GetBool() bool {
	return true
}

func (a FunctionValue) GetString() string {
	return "Function"
}
func (a FunctionValue) GetType() ValueType {
	return FunctionType
}

func (f FunctionValue) GetProp(name string) (RuntimeValue, error) {
	return NothingValue{}, fmt.Errorf("property %s does not exists", name)
}
