package runtime

import (
	environment "evie/env"
	"evie/lexer"
	"evie/lib"
	"evie/native"
	"evie/parser"
	"evie/utils"
	"evie/values"
	"fmt"
	"os"
	"strconv"
)

type Evaluator struct {
	Nodes []parser.Stmt
}

// Takes an AST and evaluates it, Node by node
func (e Evaluator) Evaluate(env *environment.Environment) *environment.Environment {
	for _, node := range e.Nodes {
		ret := e.EvaluateStmt(node, env)

		// If the return value is an ErrorValue
		if IsError(ret) {
			fmt.Println(ret.GetStr())
			os.Exit(1)
		}

	}
	return env
}

// Evaluate a single Statement node
func (e Evaluator) EvaluateStmt(node parser.Stmt, env *environment.Environment) values.RuntimeValue {

	switch n := node.(type) {
	case parser.VarDeclarationNode:
		return e.EvaluateVarDeclaration(node.(parser.VarDeclarationNode), env)
	case parser.TryCatchNode:
		return e.EvaluateTryCatchNode(node.(parser.TryCatchNode), env)
	case parser.IfStatementNode:
		return e.EvaluateIfStmt(node.(parser.IfStatementNode), env)
	case parser.ForInSatementNode:
		return e.EvaluateForInStmt(node.(parser.ForInSatementNode), env)
	case parser.FunctionDeclarationNode:
		return e.EvaluateFunctionDeclarationStmt(node.(parser.FunctionDeclarationNode), env)
	case parser.LoopStmtNode:
		return e.EvaluateLoopStmt(node.(parser.LoopStmtNode), env)
	case parser.StructMethodDeclarationNode:
		return e.EvaluateStructMethodExpression(node.(parser.StructMethodDeclarationNode), env)
	case parser.BreakNode:
		return e.EvaluateBreakNode(n, env)
	case parser.ContinueNode:
		return e.EvaluateContinueNode(n, env)
	case parser.ReturnNode:
		return e.EvaluateReturnNode(n, env)
	case parser.StructDeclarationNode:
		return e.EvaluatStructDeclarationStmt(node.(parser.StructDeclarationNode), env)
	case parser.ImportNode:
		return e.EvaluateImportNode(node.(parser.ImportNode), env)
	default: // If is not a statement, it is a expressionStmt
		return e.EvaluateExpressionStmt(node.(parser.ExpressionStmtNode), env)
	}

}

// Test for native imports
func (e Evaluator) EvaluateImportNode(node parser.ImportNode, env *environment.Environment) values.RuntimeValue {

	// Line
	line := node.Line

	// Map std libraries load methods
	libmap := lib.GetLibMap()

	// Is a native library?
	if val, ok := libmap[node.Path]; ok {
		val(env)
	} else {
		// Is a custom file

		/* Which modules were imported by the module which is being imported
		.	example:
		.
		.		File: main.ev
		.			import utils
		.
		. 	Here the main module or file try to import the utils module!
		.
		.		file: utils.ev
		.			import main
		.
		.	Here utils wants to imports main, so before we do that, we check in the import map
		.	which modules were imported by main, if utils module is there, the utils module will say:
		.	Wait! i cant import main, because main already import me!
		.	So we avoid circular imports
		.
		.
		*/
		importedByModule, _ := env.ImportMap[node.Path]

		for _, module := range importedByModule {
			if module == env.ModuleName {
				return Stop("Circular import at line " + strconv.Itoa(line) + " with module: " + node.Path)
			}
		}

		// Do not add .ev to modules ;)
		path := node.Path + ".ev"

		// Read file, parse and evaluate
		var source string = utils.ReadFile(path)

		env.ImportMap[env.ModuleName] = append(env.ImportMap[env.ModuleName], node.Path)

		tokens := lexer.Tokenize(source)
		ast := parser.NewParser(tokens).GetAST()

		// Create a new global environtment for load the module to avoid name conflicts
		parentEnv := environment.NewEnvironment(nil)
		native.SetupEnvironment(parentEnv)

		// Create new environment for the module with the parent environment
		envForModule := environment.NewEnvironment(parentEnv)
		envForModule.ModuleName = node.Path
		envForModule.ImportMap = env.ImportMap

		importEnv := Evaluator{Nodes: ast}.Evaluate(envForModule)

		// Get the created environment after evaluate the module
		// Get all the variables loaded and load into the actual environment
		// using a namespace
		env.Variables[node.Alias] = values.NamespaceValue{Value: importEnv.Variables}

	}

	return values.BooleanValue{Value: true}
}

// LOOP STATEMENT
func (e Evaluator) EvaluateLoopStmt(node parser.LoopStmtNode, env *environment.Environment) values.RuntimeValue {

	// Creating new environment
	newenv := environment.NewEnvironment(env)

	// Flag for Break statement
	keepLooping := true

	for {

		if keepLooping == false {
			return nil
		}

		// Loop through body
		for _, stmt := range node.Body {

			ret := e.EvaluateStmt(stmt, newenv)

			if ret != nil && (ret.GetType() == "ErrorValue" || ret.GetType() == "return") {
				return ret
			} else if ret != nil && ret.GetType() == "break" {
				return nil
			} else if ret != nil && ret.GetType() == "continue" {
				break
			}

		}
	}

}

// TRY CATCH
func (e Evaluator) EvaluateTryCatchNode(node parser.TryCatchNode, env *environment.Environment) values.RuntimeValue {

	body := node.Body
	catch := node.Catch

	// Creating new environment
	newenv := environment.NewEnvironment(env)

	// Loop through body
	for _, stmt := range body {

		ret := e.EvaluateStmt(stmt, newenv)

		if IsError(ret) {

			newenv.DeclareVar("error", ret)

			// CATCH BODY
			for _, cstmt := range catch {

				result := e.EvaluateStmt(cstmt, newenv)

				// Error inside the catch lol
				if result.GetType() == "ErrorValue" || result.GetType() == "return" {
					return result
				} else if result.GetType() == "break" {
					return nil
				} else if result.GetType() == "continue" {
					continue
				}
			}
			return nil
		}
	}
	return nil
}

// RETURN STMT
func (e Evaluator) EvaluateReturnNode(node parser.ReturnNode, env *environment.Environment) values.RuntimeValue {
	return values.SignalValue{
		Type:  "return",
		Value: node.Right,
		Env:   env,
	}
}

// CONTINUE
func (e Evaluator) EvaluateContinueNode(node parser.ContinueNode, env *environment.Environment) values.RuntimeValue {
	return values.SignalValue{
		Type:  "continue",
		Value: nil,
	}
}

// BREAk
func (e Evaluator) EvaluateBreakNode(node parser.BreakNode, env *environment.Environment) values.RuntimeValue {
	return values.SignalValue{
		Type:  "break",
		Value: nil,
	}
}

// FOR IN STMT
func (e Evaluator) EvaluateForInStmt(node parser.ForInSatementNode, env *environment.Environment) values.RuntimeValue {

	// Evaluate iterator expression
	iterator := e.EvaluateExpression(node.Iterator, env)

	if IsError(iterator) {
		return iterator
	}

	// Creating new environment
	newEnv := environment.NewEnvironment(env)

	// Flag for Break statement
	thereIsBreak := false

	if iterator.GetType() == "ArrayValue" {

		iterValues := iterator.(values.ArrayValue).Value

		for index, value := range iterValues {

			if thereIsBreak == true {
				break
			}

			// Load variables in env on each iteration
			newEnv.Variables[node.LocalVarName] = value

			if node.IndexVarName != "" {
				newEnv.Variables[node.IndexVarName] = values.NumberValue{Value: float64(index)}
			}

			// LOOP through for in body!
			for _, stmt := range node.Body {

				result := e.EvaluateStmt(stmt, newEnv)

				if result != nil {
					if result.GetType() == "ErrorValue" || result.GetType() == "return" {
						return result
					} else if result.GetType() == "break" {
						thereIsBreak = true
						break
					} else if result.GetType() == "continue" {
						break
					}
				}

			}
		}
	} else if iterator.GetType() == "DictionaryValue" {
		iterValues := iterator.(values.DictionaryValue).Value

		for index, value := range iterValues {

			if thereIsBreak == true {
				break
			}

			newEnv.Variables[node.IndexVarName] = values.StringValue{Value: index}
			if node.IndexVarName != "" {
				newEnv.Variables[node.LocalVarName] = value
			}

			for _, stmt := range node.Body {

				result := e.EvaluateStmt(stmt, newEnv)

				if result != nil {
					if result.GetType() == "ErrorValue" || result.GetType() == "return" {
						return result
					} else if result.GetType() == "break" {
						thereIsBreak = true
						break
					} else if result.GetType() == "continue" {
						break
					}
				}

			}

		}
	}

	return values.BooleanValue{Value: true}
}

// STRUCT DECLARATION
func (e Evaluator) EvaluatStructDeclarationStmt(node parser.StructDeclarationNode, env *environment.Environment) values.RuntimeValue {

	rtValue := values.StructValue{}

	rtValue.Properties = node.Properties
	rtValue.Methods = make(map[string]values.RuntimeValue)

	env.DeclareVar(node.Name, rtValue)

	return values.BooleanValue{Value: true}
}

// FUNCTION DECLARATION
func (e Evaluator) EvaluateFunctionDeclarationStmt(node parser.FunctionDeclarationNode, env *environment.Environment) values.RuntimeValue {

	fn := values.FunctionValue{}

	fn.Body = node.Body
	fn.Parameters = node.Parameters
	fn.Struct = ""

	env.DeclareVar(node.Name, fn)

	return values.BooleanValue{Value: true}
}

// IF STMT
func (e Evaluator) EvaluateIfStmt(node parser.IfStatementNode, env *environment.Environment) values.RuntimeValue {

	exp := node.Condition

	evaluatedExp := e.EvaluateExpression(exp, env)

	if IsError(evaluatedExp) {
		return evaluatedExp
	}

	if evaluatedExp.GetBool() == true {

		newEnv := environment.NewEnvironment(env)

		for _, stmt := range node.Body {

			result := e.EvaluateStmt(stmt, newEnv)

			if result != nil && result.GetType() == "ErrorValue" || result.GetType() == "return" || result.GetType() == "break" || result.GetType() == "continue" {
				return result
			}

		}
	} else {

		matched := false

		if node.ElseIf != nil {
			for _, elseif := range node.ElseIf {

				exp := e.EvaluateExpression(elseif.Condition, env)

				if IsError(exp) {
					return exp
				}

				if exp.GetBool() == true {

					matched = true

					newEnv := environment.NewEnvironment(env)

					for _, stmt := range elseif.Body {

						result := e.EvaluateStmt(stmt, newEnv)

						if result != nil && (result.GetType() == "ErrorValue" || result.GetType() == "return" || result.GetType() == "break" || result.GetType() == "continue") {
							return result
						}
					}
					return nil
				}
			}
		}

		if matched == false {
			newEnv := environment.NewEnvironment(env)
			for _, stmt := range node.ElseBody {

				result := e.EvaluateStmt(stmt, newEnv)
				if result != nil && (result.GetType() == "ErrorValue" || result.GetType() == "return" || result.GetType() == "break" || result.GetType() == "continue") {
					return result
				}
			}
		}
	}

	return values.BooleanValue{Value: true}
}

// Evaluate a single Expression Statement
// An Expression Statement is statement that is only formed by expressions, like assigments or binary expressions
func (e Evaluator) EvaluateExpressionStmt(node parser.ExpressionStmtNode, env *environment.Environment) values.RuntimeValue {
	return e.EvaluateExpression(node.Expression, env)
}

// Evaluate a single Variable Declaration
func (e Evaluator) EvaluateVarDeclaration(node parser.VarDeclarationNode, env *environment.Environment) values.RuntimeValue {

	// VAR NAME
	var identifier parser.IdentifierNode = node.Left

	// TODO: += -= *= ...
	// operator := node.Operator

	parsed := e.EvaluateExpression(node.Right, env)

	if IsError(parsed) {
		return parsed
	}

	env.DeclareVar(identifier.Value, parsed)

	return values.BooleanValue{Value: true}
}

// Evaluate an expression
func (e Evaluator) EvaluateExpression(node parser.Exp, env *environment.Environment) values.RuntimeValue {

	if node == nil {
		return nil
	}

	if node.ExpType() == "BinaryExpNode" {
		return e.EvaluateBinaryExpression(node.(parser.BinaryExpNode), env)
	} else if node.ExpType() == "TernaryExpNode" {
		return e.EvaluateTernaryExpression(node.(parser.TernaryExpNode), env)
	} else if node.ExpType() == "CallExpNode" {
		return e.EvaluateCallExpression(node.(parser.CallExpNode), env)
	} else if node.ExpType() == "IdentifierNode" {
		return env.GetVar(node.(parser.IdentifierNode).Value, node.(parser.IdentifierNode).Line)
	} else if node.ExpType() == "UnaryExpNode" {
		return e.EvaluateUnaryExpression(node.(parser.UnaryExpNode), env)
	} else if node.ExpType() == "NumberNode" {
		parsedNumber, _ := strconv.ParseFloat(node.(parser.NumberNode).Value, 2)
		return values.NumberValue{Value: parsedNumber}
	} else if node.ExpType() == "StringNode" {
		return values.StringValue{Value: node.(parser.StringNode).Value}
	} else if node.ExpType() == "IndexAccessExpNode" {
		return e.EvaluateIndexAccessExpression(node.(parser.IndexAccessExpNode), env)
	} else if node.ExpType() == "BooleanNode" {
		return values.BooleanValue{Value: node.(parser.BooleanNode).Value}
	} else if node.ExpType() == "AssignmentNode" {
		return e.EvaluateAssignmentExpression(node.(parser.AssignmentNode), env)
	} else if node.ExpType() == "SliceExpNode" {
		return e.EvaluateSliceExpression(node.(parser.SliceExpNode), env)
	} else if node.ExpType() == "MemberExpNode" {
		return e.EvaluateMemberExpression(node.(parser.MemberExpNode), env)
	} else if node.ExpType() == "ObjectInitExpNode" {
		return e.EvaluateObjInitializeExpression(node.(parser.ObjectInitExpNode), env)
	} else if node.ExpType() == "ArrayExpNode" {
		return e.EvaluateArrayExpression(node.(parser.ArrayExpNode), env)
	} else if node.ExpType() == "DictionaryExpNode" {
		return e.EvaluateDictionaryExpression(node.(parser.DictionaryExpNode), env)
	} else {
		return nil
	}
}

func (e Evaluator) EvaluateTernaryExpression(node parser.TernaryExpNode, env *environment.Environment) values.RuntimeValue {

	condition := e.EvaluateExpression(node.Condition, env)

	if IsError(condition) {
		return condition
	}

	if condition.GetBool() {
		return e.EvaluateExpression(node.Left, env)
	} else {
		return e.EvaluateExpression(node.Right, env)
	}
}

func (e Evaluator) EvaluateSliceExpression(node parser.SliceExpNode, env *environment.Environment) values.RuntimeValue {

	value := e.EvaluateExpression(node.Left, env)
	init := e.EvaluateExpression(node.From, env)

	if (init != nil && init.GetType() == "ErrorValue") || (value != nil && value.GetType() == "ErrorValue") {
		return init
	}

	if init == nil {
		init = values.NumberValue{Value: 0}
	}
	end := e.EvaluateExpression(node.To, env)

	if end != nil && end.GetType() == "ErrorValue" {
		return end
	}

	switch value.(type) {
	case values.ArrayValue:
		if end == nil {
			end = values.NumberValue{Value: float64(len(value.(values.ArrayValue).Value))}
		}

		fn := value.GetProp("slice")
		ret := fn.(values.NativeFunctionValue).Value([]values.RuntimeValue{init, end})
		return ret
	case values.StringValue:
		if end == nil {
			end = values.NumberValue{Value: float64(len(value.(values.StringValue).Value))}
		}

		fn := value.GetProp("slice")
		ret := fn.(values.NativeFunctionValue).Value([]values.RuntimeValue{init, end})
		return ret
	default:
		return nil
	}

}

// Declaration of a struct method
func (e Evaluator) EvaluateStructMethodExpression(node parser.StructMethodDeclarationNode, env *environment.Environment) values.RuntimeValue {

	// Struct name string
	structName := node.Struct

	// Check if struct exists and is a struct
	structLup := env.GetVar(structName, node.Line)

	if structLup == nil || structLup.GetType() != "StructValue" {
		Stop("Expected struct, got " + structLup.GetType() + " in line " + fmt.Sprint(node.Line))
	}

	// Create function value
	fn := values.FunctionValue{}

	fn.Body = node.Function.Body
	fn.Parameters = node.Function.Parameters
	fn.Struct = structName

	// Store function in struct
	env.Variables[structName].(values.StructValue).Methods[node.Function.Name] = fn

	return nil
}

// Evaluate a member expression
func (e Evaluator) EvaluateMemberExpression(node parser.MemberExpNode, env *environment.Environment) values.RuntimeValue {

	// Evaluate Base var recursively
	// baseVar.member1.member2
	// eval this first -> (baseVar.member1)
	// Left -> MemberExpNode{ Left: baseVar, Member: member1}
	varValue := e.EvaluateExpression(node.Left, env)

	if varValue.GetType() == "ErrorValue" {
		return varValue
	}

	fn := varValue.GetProp(node.Member.Value)

	if fn == nil {
		return Stop("Unknown member " + node.Member.Value + " in line " + fmt.Sprint(node.Line))
	}

	return fn

}

// Evaluate an object initialization
func (e Evaluator) EvaluateObjInitializeExpression(node parser.ObjectInitExpNode, env *environment.Environment) values.RuntimeValue {

	// Struct name string
	structName := node.Struct

	// Syntax for object initialization is the same as dictionaries
	propDict := node.Value

	val := values.ObjectValue{}

	// Lookup for struct
	val.Struct = env.GetVar(structName, node.Line).(values.StructValue)
	val.Value = make(map[string]values.RuntimeValue)

	// Evaluate expressions and set values
	// TODO: Allow only set values of properties defined in struct
	for key, exp := range propDict.Value {
		value := e.EvaluateExpression(exp, env)
		if value.GetType() == "ErrorValue" {
			return value
		}
		val.Value[key] = value
	}

	return val
}

func (e Evaluator) EvaluateDictionaryExpression(node parser.DictionaryExpNode, env *environment.Environment) values.RuntimeValue {

	dict := values.DictionaryValue{}
	dict.Value = make(map[string]values.RuntimeValue)

	for key, exp := range node.Value {
		value := e.EvaluateExpression(exp, env)
		if value.GetType() == "ErrorValue" {
			return value
		}
		dict.Value[key] = value
	}

	return dict
}

// Evaluate an index access
func (e Evaluator) EvaluateIndexAccessExpression(node parser.IndexAccessExpNode, env *environment.Environment) values.RuntimeValue {

	// Obtenemos el valor final del valor base
	identifier := e.EvaluateExpression(node.Left, env)

	if identifier.GetType() == "ErrorValue" {
		return identifier
	}

	// Obtenemos el valor final del indice
	index := e.EvaluateExpression(node.Index, env)

	if index.GetType() == "ErrorValue" {
		return identifier
	}

	// El indice puede ser numeric si es un array, o un string si se trata de un diccionario
	// En ambos casos lo trataremos como un string y si es necesario se convertira a int

	var i string = index.GetStr()

	switch val := identifier.(type) {
	case values.ArrayValue:
		iToInt, _ := strconv.Atoi(i)
		return val.Value[iToInt]
	case values.StringValue:
		iToInt, _ := strconv.Atoi(i)
		return values.StringValue{Value: string(val.Value[iToInt])}
	case values.DictionaryValue:
		item, exists := val.Value[i]

		if !exists {
			return Stop("Undefined key '" + i + "' in line " + fmt.Sprint(node.Line))
		}

		return item

	default:
		return Stop("At line " + fmt.Sprint(node.Line) + ":\nOnly arrays and dictionaries can be accessed by index")
	}

}

func (e Evaluator) EvaluateArrayExpression(node parser.ArrayExpNode, env *environment.Environment) values.RuntimeValue {

	rtvalue := values.ArrayValue{}
	rtvalue.Value = make([]values.RuntimeValue, 0)

	for _, exp := range node.Value {
		rtvalue.Value = append(rtvalue.Value, e.EvaluateExpression(exp, env))
	}

	return rtvalue

}

func (e Evaluator) EvaluateCallExpression(node parser.CallExpNode, env *environment.Environment) values.RuntimeValue {

	evaluatedArgs := []values.RuntimeValue{}

	for _, arg := range node.Args {
		value := e.EvaluateExpression(arg, env)
		if value.GetType() == "ErrorValue" {
			return value
		}
		evaluatedArgs = append(evaluatedArgs, value)

	}

	calle := e.EvaluateExpression(node.Name, env)

	if calle.GetType() == "ErrorValue" {
		return calle
	}

	switch fn := calle.(type) {

	case values.ErrorValue:
		return fn
	case values.NativeFunctionValue:
		return fn.Value(evaluatedArgs)
	case values.FunctionValue:

		// Create new environment
		newEnv := environment.NewEnvironment(env)

		// Set this
		if fn.Struct != "" {
			newEnv.DeclareVar("this", fn.StructObjRef)
		}

		for index, param := range fn.Parameters {
			value := e.EvaluateExpression(node.Args[index], env)
			if value.GetType() == "ErrorValue" {
				return value
			}
			newEnv.DeclareVar(param, value)
		}

		var result values.RuntimeValue

		for _, stmt := range fn.Body {
			result = e.EvaluateStmt(stmt, newEnv)

			if result != nil && result.GetType() == "ErrorValue" {
				return result
			}

			signal, isSignal := result.(values.SignalValue)

			if isSignal && signal.Type == "return" {
				exp := e.EvaluateExpression(signal.Value, signal.Env.(*environment.Environment))

				if exp.GetType() == "ErrorValue" {
					return exp
				}

				return exp
			}

		}

	}

	return values.BooleanValue{Value: true}
}

// Creates an access props chain
// examples:
// exp -> a[b][c] -> returns ->  ["a", "b", "c"]
// exp -> a[2][3] -> returns ->  ["a", "2", "3"]
// exp -> a[2][x] -> returns ->  ["a", "2", "x"]
func (e *Evaluator) ResolveIndexAccessChain(node parser.IndexAccessExpNode, env *environment.Environment) []string {

	indexes := []string{}

	// the last index is the last of the chain
	var lastIndex string

	// The last index can be literal or an expression
	switch index := node.Index.(type) {

	case parser.NumberNode:
		lastIndex = index.Value

	default:
		// if one of the index is not a literal value, like a[x], need to eval the expression
		numValue := e.EvaluateExpression(index, env)
		if numValue.GetType() == "ErrorValue" {
			return []string{}
		}
		lastIndex = numValue.GetStr()
	}

	indexes = append(indexes, lastIndex)

	// if the base Var is another index access, resolve the chain recursively
	if node.Left.ExpType() == "IndexAccessExpNode" {
		indexes = append(e.ResolveIndexAccessChain(node.Left.(parser.IndexAccessExpNode), env), indexes...)
	} else if node.Left.ExpType() == "IdentifierNode" {
		valName := node.Left.(parser.IdentifierNode).Value
		indexes = append([]string{valName}, indexes...)
	}

	return indexes
}

// Evaluate an Assignment Expression
func (e Evaluator) EvaluateAssignmentExpression(node parser.AssignmentNode, env *environment.Environment) values.RuntimeValue {
	// Left can be an expression, like a index access, or member, etc.
	// We will make the env function to resolve this.
	left := node.Left

	// Evaluate the expression of the right side
	right := e.EvaluateExpression(node.Right, env)

	if right.GetType() == "ErrorValue" {
		return right
	}

	if left.ExpType() == "IndexAccessExpNode" {
		chain := e.ResolveIndexAccessChain(left.(parser.IndexAccessExpNode), env)
		env.ModifyIndexValue(left.(parser.IndexAccessExpNode), right, chain)
		return right
	} else {
		env.SetVar(left, right)
	}

	return right
}

func (e Evaluator) EvaluateBinaryExpression(node parser.BinaryExpNode, env *environment.Environment) values.RuntimeValue {

	left := e.EvaluateExpression(node.Left, env)
	right := e.EvaluateExpression(node.Right, env)

	if left.GetType() == "ErrorValue" || right.GetType() == "ErrorValue" {
		return left
	}

	op := node.Operator

	types := left.GetType()

	if types != right.GetType() {
		return Stop("Type mismatch in line: " + fmt.Sprint(node.Line) + ". " + types + " and " + right.GetType())
	}

	if op == "+" {
		if types == "StringValue" {
			return values.StringValue{Value: left.GetStr() + right.GetStr()}
		} else if types == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() + right.GetNumber()}
		}
	} else if op == "-" {
		if types == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() - right.GetNumber()}
		} else {
			return Stop("Cant subtract with strings in line: " + fmt.Sprint(node.Line))
		}
	} else if op == "*" {
		if types == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() * right.GetNumber()}
		} else {
			Stop("Cant multiply with strings in line: " + fmt.Sprint(node.Line))
		}
	} else if op == "/" {
		if types == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() / right.GetNumber()}
		} else {
			Stop("Cant divide with strings in line: " + fmt.Sprint(node.Line))
		}
	} else if op == "==" {
		if types == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() == right.GetNumber()}
		} else if types == "StringValue" {
			return values.BooleanValue{Value: left.GetStr() == right.GetStr()}
		} else if types == "BooleanValue" {
			return values.BooleanValue{Value: left.GetBool() == right.GetBool()}
		}
	} else if op == ">" {
		if types == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() > right.GetNumber()}
		}
	} else if op == "<" {
		if types == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() < right.GetNumber()}
		}
	} else if op == "<=" {
		if types == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() <= right.GetNumber()}
		}
	} else if op == ">=" {
		if types == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() >= right.GetNumber()}
		}
	} else if op == "and" {
		if types == "BooleanValue" {
			return values.BooleanValue{Value: left.GetBool() && right.GetBool()}
		}
	} else if op == "or" {
		if types == "BooleanValue" {
			return values.BooleanValue{Value: left.GetBool() || right.GetBool()}
		}
	}

	return nil

}

func (e Evaluator) EvaluateUnaryExpression(node parser.UnaryExpNode, env *environment.Environment) values.RuntimeValue {

	exp := e.EvaluateExpression(node.Right, env)

	if node.Operator == "-" && exp.GetType() == "NumberValue" {
		return values.NumberValue{Value: -exp.GetNumber()}
	} else {
		return nil
	}
}

func Stop(msg string) values.ErrorValue {
	return values.ErrorValue{Value: msg}
}

func IsError(val values.RuntimeValue) bool {
	return (val != nil && val.GetType() == "ErrorValue")

}
