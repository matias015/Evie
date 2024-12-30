package environment

import "evie/values"

// Creates a instance of an Environment
func NewEnvironment() *Environment {
	env := Environment{
		Variables:   make([]map[string]values.RuntimeValue, 0),
		ImportChain: make(map[string]bool),
	}

	return &env
}
