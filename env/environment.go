package environment

import (
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

	// Variables is a map of variable names to their values
	Variables []map[string]values.RuntimeValue

	//
	ScopeCount int

	// keep tracking of imports flow to avoid circular imports
	ImportChain map[string]bool

	// Module Name
	ModuleName string
}

func (e *Environment) PushScope() {
	e.Variables = append(e.Variables, e.GetFromPool())
	e.ScopeCount++
}

func (e *Environment) ExitScope() {
	last := e.GetCurrentScope()
	e.PutToPool(last)
	e.Variables = e.Variables[:e.GetCurrentScopeCount()]
	e.ScopeCount--
}

// Declares a variable
func (e *Environment) DeclareVar(name string, value values.RuntimeValue) error {
	currentScope := e.GetCurrentScope()

	if _, exists := currentScope[name]; exists {
		return fmt.Errorf("variable '%s' already declared", name)
	}
	currentScope[name] = value
	return nil
}

// Returns the value of a variable
func (e *Environment) GetVar(name string) (values.RuntimeValue, error) {

	for i := e.GetCurrentScopeCount(); i >= 0; i-- {
		if val, exists := e.Variables[i][name]; exists {
			return val, nil
		}
	}
	return values.NothingValue{}, fmt.Errorf("variable '%s' not found", name)
}

func (e *Environment) ForceDeclare(name string, value values.RuntimeValue) {
	e.Variables[e.GetCurrentScopeCount()][name] = value
}

// Assigns a value to a variable
func (e *Environment) SetVar(name string, value values.RuntimeValue) error {

	for i := e.GetCurrentScopeCount(); i >= 0; i-- { // Recorre desde el último alcance
		scope := e.Variables[i]
		if _, exists := scope[name]; exists {
			scope[name] = value // Actualiza el valor
			return nil          // Operación exitosa
		}
	}

	return fmt.Errorf("variable '%s' not found", name)

}

func (e *Environment) GetFromPool() map[string]values.RuntimeValue {
	return mapPool.Get().(map[string]values.RuntimeValue)
}

func (e *Environment) PutToPool(last map[string]values.RuntimeValue) {
	for k := range last {
		delete(last, k)
	}

	mapPool.Put(last)
}
func (e *Environment) GetCurrentScope() map[string]values.RuntimeValue {
	return e.Variables[e.ScopeCount]
}

func (e *Environment) GetCurrentScopeCount() int {
	return e.ScopeCount
}
