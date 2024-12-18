package environment

import (
	"evie/parser"
	"evie/values"
	"fmt"
	"strconv"
)

type Environment struct {
	// Parent evironment, if a variable is not found in the current environment
	// it will be searched in the parent
	Parent *Environment

	// Variables is a map of variable names to their values
	Variables map[string]values.RuntimeValue

	// Imported Files
	ImportMap map[string][]string

	// Module Name
	ModuleName string
}

// Declares a variable
func (e *Environment) DeclareVar(name string, value values.RuntimeValue) {
	e.Variables[name] = value
}

// Returns the value of a variable
func (e *Environment) GetVar(name string, line int) (values.RuntimeValue, error) {

	// Look for the variable in the current environment
	if value, ok := e.Variables[name]; ok {
		return value, nil
	}

	// If the variable is not found in the current environment
	if e.Parent != nil {
		return e.Parent.GetVar(name, line)
	}

	return nil, fmt.Errorf("Undefined variable '%s'", name)
}

// Assigns a value to a variable
func (e *Environment) SetVar(name parser.Exp, value values.RuntimeValue) {

	switch left := name.(type) {
	case parser.IdentifierNode: // it is simple asignment like a = Exp
		_, exists := e.Variables[left.Value]
		if !exists {
			if e.Parent != nil {
				e.Parent.SetVar(left, value)
			} else {
				Stop("Undefined variable '"+left.Value+"'", left.Line, e.ModuleName)
			}
		}
		e.Variables[left.Value] = value
	// case parser.IndexAccessExpNode: // it is index assignment like a[i] = Exp or a[b][c] = Exp
	//  -> This is made by the evaluator
	case parser.MemberExpNode: // it is index assignment like a[i] = Exp or a[b][c] = Exp
		e.ModifyByMember(left, value)
	}
}

// Assigns a value to a variable by member
// ex a.b = 1 or a.b.c = 1
func (e *Environment) ModifyByMember(left parser.MemberExpNode, value values.RuntimeValue) {
	// Create a access props chain
	chain := e.ResolveMemberAccessChain(left)

	// Assign the value
	e.ModifyMemberValue(left, value, chain)
}

func (e *Environment) ModifyMemberValue(left parser.MemberExpNode, value values.RuntimeValue, chain []string) {
	// Get the initial element of the chain, which is the base variable
	varName := chain[0]

	// Get the value of the base variable
	// looping through the its properties by following the chain will be turned into the final value
	endValue, exists := e.Variables[varName]
	// fmt.Println(chain)

	if !exists {
		if e.Parent != nil {
			e.Parent.ModifyMemberValue(left, value, chain)
			return
		} else {
			Stop("Undefined variable '"+varName+"'", left.Line, e.ModuleName)
		}
	}

	// only use the middle of the chain, if chain is ["a", "b", "2", "c"] then only use ["b", "2"]
	// Because the first element of the chain is the base variable
	// and the last will be used after ending the loop to make the final assigment
	middleChain := chain[1 : len(chain)-1]

	for _, el := range middleChain {

		// In each iteration, access to the next property of the chain in
		// the value should be an array
		switch val := endValue.(type) {
		case values.ObjectValue: // each prop of chain is a str so convert to int
			endValue = val.Value[el] // access to the next property of the chain in the array
		default:
			Stop("Only objects can be accessed by dot notation", left.Line, e.ModuleName)
		}
	}

	// Assign the value
	switch val := endValue.(type) {
	case values.ObjectValue:
		// Get the last element of the chain and convert it to int
		lastIndex := chain[len(chain)-1]
		// Assign the value
		val.Value[lastIndex] = value
	default:
		Stop("Only objects can be accessed by dot notation", left.Line, e.ModuleName)
	}
}

func (e *Environment) ModifyIndexValue(left parser.IndexAccessExpNode, value values.RuntimeValue, chain []string) values.RuntimeValue {
	// Get the initial element of the chain, which is the base variable
	varName := chain[0]

	// Get the value of the base variable
	// looping through the its properties by following the chain will be turned into the final value
	endValue, exists := e.Variables[varName]

	if !exists {
		if e.Parent != nil {
			return e.Parent.ModifyIndexValue(left, value, chain)
		} else {
			return Stop("Undefined variable '"+varName+"'", left.Line, e.ModuleName)
		}
	}

	// only use the middle of the chain, if chain is ["a", "b", "2", "c"] then only use ["b", "2"]
	// Because the first element of the chain is the base variable
	// and the last will be used after ending the loop to make the final assigment
	middleChain := chain[1 : len(chain)-1]

	for _, el := range middleChain {

		// In each iteration, access to the next property of the chain in
		// the value should be an array
		switch val := endValue.(type) {
		case values.ArrayValue:
			elToInt, _ := strconv.Atoi(el) // each prop of chain is a str so convert to int
			endValue = val.Value[elToInt]  // access to the next property of the chain in the array
		case values.DictionaryValue:
			endValue = val.Value[el]
		case values.StringValue:
			endValue = val
		default:
			return Stop("Only arrays and dictionaries can be accessed by index", left.Line, e.ModuleName)
		}
	}

	// Assign the value
	switch val := endValue.(type) {
	case values.ArrayValue:
		// Get the last element of the chain and convert it to int
		lastIndex, _ := strconv.Atoi(chain[len(chain)-1])
		// Assign the value
		val.Value[lastIndex] = value
	case values.DictionaryValue:
		// Get the last element of the chain and convert it to int
		lastIndex := chain[len(chain)-1]
		// Assign the value
		val.Value[lastIndex] = value

	default:
		return Stop("Only arrays and dictionaries can be assigned by index", left.Line, e.ModuleName)
	}
	return values.BooleanValue{Value: true}
}

// Creates an access props chain
// examples:
// exp -> a[b][c] -> returns ->  ["a", "b", "c"]
// exp -> a[2][3] -> returns ->  ["a", "2", "3"]
// exp -> a[2][x] -> returns ->  ["a", "2", "x"]
func (e *Environment) ResolveMemberAccessChain(node parser.MemberExpNode) []string {

	indexes := []string{}

	indexes = append(indexes, node.Member.Value)

	// if the base Var is another index access, resolve the chain recursively
	if node.Left.ExpType() == "MemberExpNode" {
		indexes = append(e.ResolveMemberAccessChain(node.Left.(parser.MemberExpNode)), indexes...)
	} else if node.Left.ExpType() == "IdentifierNode" {
		valName := node.Left.(parser.IdentifierNode).Value
		indexes = append([]string{valName}, indexes...)
	}

	return indexes
}

func Stop(msg string, line int, mod string) values.ErrorValue {
	output := "Error in line " + fmt.Sprint(line) + " at module " + mod + ":\n" + msg + "\n"
	return values.ErrorValue{Value: output}
}
