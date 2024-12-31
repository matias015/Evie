package environment

import (
	"evie/parser"
	"evie/values"
	"fmt"
	"sync"
)

var mapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]values.RuntimeValue)
	},
}

type Environment struct {
	// Parent evironment, if a variable is not found in the current environment
	// it will be searched in the parent
	Parent *Environment

	// Variables is a map of variable names to their values
	Variables []map[string]values.RuntimeValue

	//
	ImportChain map[string]bool

	// Module Name
	ModuleName string
}

func (e *Environment) PushScope() {
	val := mapPool.Get().(map[string]values.RuntimeValue)
	e.Variables = append(e.Variables, val)
}

func (e *Environment) ExitScope() {
	len := len(e.Variables) - 1
	last := e.Variables[len]
	for k := range last {
		delete(last, k)
	}

	mapPool.Put(last)
	e.Variables = e.Variables[:len]
}

// Declares a variable
func (e *Environment) DeclareVar(name string, value values.RuntimeValue) error {
	currentScope := e.Variables[len(e.Variables)-1]
	if _, exists := currentScope[name]; exists {
		return fmt.Errorf("variable '%s' already declared", name)
	}
	currentScope[name] = value
	return nil
}

// Returns the value of a variable
func (e *Environment) GetVar(name string, line int) (values.RuntimeValue, error) {

	for i := len(e.Variables) - 1; i >= 0; i-- {
		if val, exists := e.Variables[i][name]; exists {
			return val, nil
		}
	}
	return values.NothingValue{}, fmt.Errorf("variable '%s' not found", name)
}

func (e *Environment) ForceDeclare(name string, value values.RuntimeValue) {
	e.Variables[len(e.Variables)-1][name] = value
}

// Assigns a value to a variable
func (e *Environment) SetVar(name parser.Exp, value values.RuntimeValue) error {

	switch left := name.(type) {
	case parser.IdentifierNode: // it is simple asignment like a = Exp
		for i := len(e.Variables) - 1; i >= 0; i-- { // Recorre desde el último alcance
			scope := e.Variables[i]
			if _, exists := scope[left.Value]; exists {
				scope[left.Value] = value // Actualiza el valor
				return nil                // Operación exitosa
			}
		}
		return fmt.Errorf("variable '%s' not found", left.Value)
	default:
		return fmt.Errorf("Invalid assignment")
	}
}
