package environment

import "evie/values"

// Creates a instance of an Environment
func NewEnvironment(parent *Environment) *Environment {
	env := Environment{
		Parent:    parent,
		Variables: make(map[string]values.RuntimeValue),
		ImportMap: make(map[string][]string),
	}

	if parent != nil {
		env.ImportMap = parent.ImportMap
		env.ModuleName = parent.ModuleName
	}

	return &env
}
