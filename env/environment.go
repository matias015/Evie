package environment

import (
	"evie/parser"
	"evie/values"
	"fmt"
	"strconv"
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

	// Imported Files
	ImportMap map[string][]string

	// Module Name
	ModuleName string
}

func (e *Environment) PushScope() {
	myMap := mapPool.Get().(map[string]values.RuntimeValue)
	e.Variables = append(e.Variables, myMap)
}

func (e *Environment) ExitScope() {
	for k := range e.Variables[len(e.Variables)-1] {
		delete(e.Variables[len(e.Variables)-1], k)
	}
	mapPool.Put(e.Variables[len(e.Variables)-1])
	e.Variables = e.Variables[:len(e.Variables)-1]
}

// Declares a variable
func (e *Environment) DeclareVar(name string, value values.RuntimeValue) (bool, error) {
	currentScope := e.Variables[len(e.Variables)-1]
	if _, exists := currentScope[name]; exists {
		return false, fmt.Errorf("variable '%s' already declared", name)
	}
	currentScope[name] = value
	return true, nil
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

func (e *Environment) ForceDeclare(name string, value values.RuntimeValue) (bool, error) {
	e.Variables[len(e.Variables)-1][name] = value
	return true, nil
}

// Assigns a value to a variable
func (e *Environment) SetVar(name parser.Exp, value values.RuntimeValue) (bool, error) {

	switch left := name.(type) {
	case parser.IdentifierNode: // it is simple asignment like a = Exp
		for i := len(e.Variables) - 1; i >= 0; i-- { // Recorre desde el último alcance
			scope := e.Variables[i]
			if _, exists := scope[left.Value]; exists {
				scope[left.Value] = value // Actualiza el valor
				return true, nil          // Operación exitosa
			}
		}
		return false, fmt.Errorf("variable '%s' not found", left.Value)
	// case parser.IndexAccessExpNode: // it is index assignment like a[i] = Exp or a[b][c] = Exp
	//  -> This is made by the evaluator
	case parser.MemberExpNode: // it is index assignment like a[i] = Exp or a[b][c] = Exp
		return e.ModifyByMember(left, value)
	default:
		return false, fmt.Errorf("Invalid assignment")
	}
}

// Assigns a value to a variable by member
// ex a.b = 1 or a.b.c = 1
func (e *Environment) ModifyByMember(left parser.MemberExpNode, value values.RuntimeValue) (bool, error) {
	// Create a access props chain
	chain := e.ResolveMemberAccessChain(left)

	// Assign the value
	return e.ModifyMemberValue(left, value, chain)
}

func (e *Environment) ModifyMemberValue(left parser.MemberExpNode, value values.RuntimeValue, chain []string) (bool, error) {
	// Get the initial element of the chain, which is the base variable
	varName := chain[0]

	// Get the value of the base variable
	// looping through the its properties by following the chain will be turned into the final value
	endValue, err := e.GetVar(varName, left.Line)
	// fmt.Println(chain)

	if err != nil {
		return false, fmt.Errorf("Undefined variable '" + varName + "'")
	}

	// only use the middle of the chain, if chain is ["a", "b", "2", "c"] then only use ["b", "2"]
	// Because the first element of the chain is the base variable
	// and the last will be used after ending the loop to make the final assigment
	middleChain := chain[1 : len(chain)-1]

	for _, el := range middleChain {

		// In each iteration, access to the next property of the chain in
		// the value should be an array
		switch endValue.GetType() {
		case values.ObjectType: // each prop of chain is a str so convert to int
			endValue = endValue.(values.ObjectValue).Value[el] // access to the next property of the chain in the array
		default:
			return false, fmt.Errorf("Only objects can be accessed by dot notation")
		}
	}

	// Assign the value
	switch endValue.GetType() {
	case values.ObjectType:

		// Get the last element of the chain and convert it to int
		lastIndex := chain[len(chain)-1]

		_, propExists := endValue.(*values.ObjectValue).Value[lastIndex]

		if !propExists {
			return false, fmt.Errorf("Undefined property '" + chain[len(chain)-1] + "' in object")
		}

		// Assign the value
		endValue.(*values.ObjectValue).Value[lastIndex] = value

		return true, nil
	default:
		return false, fmt.Errorf("Only objects can be accessed by dot notation")
	}
}

func (e *Environment) ModifyIndexValue(left parser.IndexAccessExpNode, value values.RuntimeValue, chain []string) (bool, error) {
	// Get the initial element of the chain, which is the base variable
	varName := chain[0]

	// Get the value of the base variable
	// looping through the its properties by following the chain will be turned into the final value
	endValue, err := e.GetVar(varName, left.Line)

	if err != nil {
		return false, fmt.Errorf("Undefined variable '" + varName + "'")
	}

	// only use the middle of the chain, if chain is ["a", "b", "2", "c"] then only use ["b", "2"]
	// Because the first element of the chain is the base variable
	// and the last will be used after ending the loop to make the final assigment
	middleChain := chain[1 : len(chain)-1]

	for _, el := range middleChain {

		// In each iteration, access to the next property of the chain in
		// the value should be an array
		switch endValue.GetType() {
		case values.ArrayType:
			elToInt, _ := strconv.Atoi(el)                          // each prop of chain is a str so convert to int
			endValue = endValue.(*values.ArrayValue).Value[elToInt] // access to the next property of the chain in the array
		case values.DictionaryType:
			endValue = endValue.(*values.DictionaryValue).Value[el]
		default:
			return false, fmt.Errorf("Only arrays and dictionaries can be modified by index access")
		}
	}

	// Assign the value
	switch endValue.GetType() {
	case values.ArrayType:
		// Get the last element of the chain and convert it to int
		lastIndex, _ := strconv.Atoi(chain[len(chain)-1])
		// Assign the value
		if lastIndex >= len(endValue.(*values.ArrayValue).Value) {
			return false, fmt.Errorf("Index " + chain[len(chain)-1] + " out of range")
		}
		endValue.(*values.ArrayValue).Value[lastIndex] = value
	case values.DictionaryType:
		// Get the last element of the chain and convert it to int
		lastIndex := chain[len(chain)-1]

		_, propExists := endValue.(*values.DictionaryValue).Value[lastIndex]

		if !propExists {
			return false, fmt.Errorf("Undefined property '" + chain[len(chain)-1] + "' in dictionary")
		}

		// Assign the value
		endValue.(*values.DictionaryValue).Value[lastIndex] = value
	default:
		return false, fmt.Errorf("Only arrays and dictionaries can be modified by index access")
	}
	return true, nil
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
	if node.Left.ExpType() == parser.NodeMemberExp {
		indexes = append(e.ResolveMemberAccessChain(node.Left.(parser.MemberExpNode)), indexes...)
	} else if node.Left.ExpType() == parser.NodeIdentifier {
		valName := node.Left.(parser.IdentifierNode).Value
		indexes = append([]string{valName}, indexes...)
	}

	return indexes
}

func Stop(msg string, line int, mod string) values.RuntimeValue {
	output := "Error in line " + fmt.Sprint(line) + " at module " + mod + ":\n" + msg + "\n"
	return values.ErrorValue{Value: output}
}
