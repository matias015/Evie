package core

type Evaluator interface {
	ExecuteCallback(interface{}, interface{}, []interface{}) interface{}
}
