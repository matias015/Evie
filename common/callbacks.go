package common

type Evaluator interface {
	ExecuteCallback(interface{}, []interface{}) interface{}
}
