package environment

import (
	"evie/profil"
	"evie/values"
	"fmt"
	"sync"
)

var timer *profil.Timer = profil.ObtenerInstancia()

var mapPool = sync.Pool{
	New: func() interface{} {
		return Environment{}
	},
}

type Environment struct {

	// Variables is a map of variable names to their values
	Variables map[string]values.RuntimeValue

	Parent *Environment

	// keep tracking of imports flow to avoid circular imports
	ImportChain map[string]bool

	// Module Name
	ModuleName string
}

// Declares a variable
func (e *Environment) DeclareVar(name string, value values.RuntimeValue) error {
	//init := timer.Init()
	currentScope := e.GetCurrentScope()

	if _, exists := currentScope[name]; exists {
		return fmt.Errorf("variable '%s' already declared", name)
	}
	currentScope[name] = value
	// timer.Add("env_declare", init)
	return nil
}

// Returns the value of a variable
func (e *Environment) GetVar(name string) (values.RuntimeValue, error) {

	//init := timer.Init()

	if val, exists := e.Variables[name]; exists {
		return val, nil
	}

	if e.Parent != nil {
		return e.Parent.GetVar(name)
	}

	// timer.Add("env_lookup", init)
	return values.NothingValue{}, fmt.Errorf("variable '%s' not found", name)
}

func (e *Environment) CheckVarExists(name string) bool {

	_, exists := e.Variables[name]

	if exists {
		return true
	}

	if e.Parent != nil {
		return e.Parent.CheckVarExists(name)
	}

	return false
}

func (e *Environment) ForceDeclare(name string, value values.RuntimeValue) {
	// //init := timer.Init()
	e.Variables[name] = value
	// timer.Add("env_declare", init)
}

// Assigns a value to a variable
func (e *Environment) SetVar(name string, value values.RuntimeValue) error {

	_, ex := e.Variables[name]

	if ex {
		e.Variables[name] = value
		return nil
	} else {
		if e.Parent != nil {
			return e.Parent.SetVar(name, value)
		}
	}

	return fmt.Errorf("variable '%s' not found", name)

}

func (e *Environment) GetCurrentScope() map[string]values.RuntimeValue {
	return e.Variables
}

func NewScopeEnv(parent *Environment, size int) *Environment {
	return &Environment{
		Parent:      parent,
		Variables:   make(map[string]values.RuntimeValue, size),
		ImportChain: parent.ImportChain,
		ModuleName:  parent.ModuleName,
	}
}
