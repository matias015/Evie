package environment

import (
	"evie/values"
	"fmt"
)

type Environment struct {

	// Variables is a map of variable names to their values
	Variables map[string]values.RuntimeValue

	// keep tracking of imports flow to avoid circular imports
	ImportChain map[string]bool

	// Module Name
	ModuleName string

	Parent *Environment
}

// Declares a variable
func (env *Environment) DeclareVar(name string, value values.RuntimeValue) error {

	if env.ExistsInActualEnv(name) {
		return fmt.Errorf("variable '%s' already declared", name)
	}

	env.ForceDeclare(name, value)
	return nil
}

func (env Environment) ExistsInActualEnv(n string) bool {
	_, ex := env.Variables[n]
	return ex
}

// Returns the value of a variable
func (env *Environment) GetVar(name string) (values.RuntimeValue, error) {

	val, exists := env.Variables[name]

	if exists {
		return val, nil
	}

	if env.Parent != nil {
		return env.Parent.GetVar(name)
	}

	return values.NothingValue{}, fmt.Errorf("variable '%s' not found", name)
}

func (env Environment) CheckVarExists(name string) bool {

	_, exists := env.Variables[name]

	if exists {
		return true
	}

	if env.Parent != nil {
		return env.Parent.CheckVarExists(name)
	}

	return false
}

func (env *Environment) ForceDeclare(name string, value values.RuntimeValue) {
	env.Variables[name] = value
}

// Assigns a value to a variable
func (env *Environment) SetVar(name string, value values.RuntimeValue) error {

	if env.ExistsInActualEnv(name) {
		env.ForceDeclare(name, value)
		return nil
	} else {
		if env.Parent != nil {
			return env.Parent.SetVar(name, value)
		}
	}

	return fmt.Errorf("variable '%s' not found", name)
}

func NewScopeEnv(parent *Environment, size int) *Environment {
	return &Environment{
		Parent:      parent,
		Variables:   make(map[string]values.RuntimeValue, size),
		ImportChain: parent.ImportChain,
		ModuleName:  parent.ModuleName,
	}
}
