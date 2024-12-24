package evruntime

import (
	"errors"
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
	Nodes     []parser.Stmt
	CallStack []string
}

// Takes an AST and evaluates it, Node by node
func (e *Evaluator) Evaluate(env *environment.Environment) *environment.Environment {

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
func (e *Evaluator) EvaluateStmt(n parser.Stmt, env *environment.Environment) values.RuntimeValue {
	switch n.StmtType() {
	case parser.NodeExpStmt:
		return e.EvaluateExpressionStmt(n.(parser.ExpressionStmtNode), env)
	case parser.NodeVarDeclaration:
		return e.EvaluateVarDeclaration(n.(parser.VarDeclarationNode), env)
	case parser.NodeIfStatement:
		return e.EvaluateIfStmt(n.(parser.IfStatementNode), env)
	case parser.NodeForInStatement:
		return e.EvaluateForInStmt(n.(parser.ForInSatementNode), env)
	case parser.NodeFunctionDeclaration:
		return e.EvaluateFunctionDeclarationStmt(n.(parser.FunctionDeclarationNode), env)
	case parser.NodeLoopStatement:
		return e.EvaluateLoopStmt(n.(parser.LoopStmtNode), env)
	case parser.NodeStructMethodDeclaration:
		return e.EvaluateStructMethodExpression(n.(parser.StructMethodDeclarationNode), env)
	case parser.NodeBreakStatement:
		return e.EvaluateBreakNode(n.(parser.BreakNode), env)
	case parser.NodeContinueStatement:
		return e.EvaluateContinueNode(n.(parser.ContinueNode), env)
	case parser.NodeReturnStatement:
		return e.EvaluateReturnNode(n.(parser.ReturnNode), env)
	case parser.NodeTryCatch:
		return e.EvaluateTryCatchNode(n.(parser.TryCatchNode), env)
	case parser.NodeStructDeclaration:
		return e.EvaluatStructDeclarationStmt(n.(parser.StructDeclarationNode), env)
	case parser.NodeImportStatement:
		return e.EvaluateImportNode(n.(parser.ImportNode), env)
	default: // If is not a statement, it is a expressionStmt
		return e.EvaluateExpressionStmt(n.(parser.ExpressionStmtNode), env)
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
		.	which modules were imported by main, if utils module is there, the utils module will say circular
		.
		.
		*/
		// TODO: Detect circular imports in import chain
		// ex: a imports b, b imports c, and c imports a
		importedByModule, _ := env.ImportMap[node.Path]

		for _, module := range importedByModule {
			if module == env.ModuleName {
				return e.Panic("Circular import with module: "+node.Path, line, env)
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
		parentEnv := environment.NewEnvironment()
		native.SetupEnvironment(parentEnv)
		parentEnv.ModuleName = node.Path

		// Create new environment for the module with the parent environment
		envForModule := environment.NewEnvironment()
		envForModule.ImportMap = env.ImportMap

		eval := Evaluator{Nodes: ast}
		importEnv := eval.Evaluate(envForModule)
		// Get the created environment after evaluate the module
		// Get all the variables loaded and load into the actual environment
		// using a namespace
		env.DeclareVar(node.Alias, values.NamespaceValue{Value: importEnv.Variables[0]})

	}

	return values.BooleanValue{Value: true}
}

// LOOP STATEMENT
func (e Evaluator) EvaluateLoopStmt(node parser.LoopStmtNode, env *environment.Environment) values.RuntimeValue {

	// Creating new environment
	env.PushScope()

	for {

		// Loop through body
		for _, stmt := range node.Body {

			ret := e.EvaluateStmt(stmt, env)

			if ret != nil && (ret.GetType() == "ErrorValue" || ret.GetType() == "return") {
				env.ExitScope()
				return ret
			} else if ret != nil && ret.GetType() == "break" {
				env.ExitScope()
				return nil
			} else if ret != nil && ret.GetType() == "continue" {
				break
			}

		}
	}
	env.ExitScope()
	return nil
}

// TRY CATCH
func (e Evaluator) EvaluateTryCatchNode(node parser.TryCatchNode, env *environment.Environment) values.RuntimeValue {

	body := node.Body
	catch := node.Catch

	// Creating new environment
	env.PushScope()

	// Loop through body
	for _, stmt := range body {

		ret := e.EvaluateStmt(stmt, env)

		if ret != nil && ret.GetType() == "return" {
			if node.Finally != nil {
				e.EvaluateFinallyBlock(node.Finally, env)
			}
			env.ExitScope()
			return ret
		}

		if IsError(ret) {

			env.DeclareVar("error", values.CapturedErrorValue{Value: ret.(values.ErrorValue)})

			// CATCH BODY
			for _, cstmt := range catch {

				result := e.EvaluateStmt(cstmt, env)

				// Error inside the catch lol
				if result.GetType() == "ErrorValue" || result.GetType() == "return" {
					if node.Finally != nil {
						e.EvaluateFinallyBlock(node.Finally, env)
					}
					env.ExitScope()
					return result
				} else if result.GetType() == "break" {
					env.ExitScope()
					return nil
				} else if result.GetType() == "continue" {
					continue
				}
			}

			if node.Finally != nil {
				e.EvaluateFinallyBlock(node.Finally, env)
			}
			env.ExitScope()
			return nil
		}
	}
	if node.Finally != nil {
		e.EvaluateFinallyBlock(node.Finally, env)
	}
	env.ExitScope()
	return nil
}

func (e Evaluator) EvaluateFinallyBlock(stmt []parser.Stmt, env *environment.Environment) values.RuntimeValue {

	fmt.Println("Ejecutando finally")

	return values.BooleanValue{Value: true}

}

// RETURN STMT
func (e Evaluator) EvaluateReturnNode(node parser.ReturnNode, env *environment.Environment) values.RuntimeValue {
	ret := values.SignalValue{
		Type: "return",
		Exp:  node.Right,
		Env:  env,
	}

	eval := e.EvaluateExpression(ret.Exp, env)
	// litter.Dump(eval)
	if eval.GetType() == "ErrorValue" {
		return eval
	}

	ret.Value = eval

	return ret
}

// CONTINUE
func (e Evaluator) EvaluateContinueNode(node parser.ContinueNode, env *environment.Environment) values.RuntimeValue {
	return values.SignalValue{
		Type: "continue",
		Exp:  nil,
	}
}

// BREAk
func (e Evaluator) EvaluateBreakNode(node parser.BreakNode, env *environment.Environment) values.RuntimeValue {
	return values.SignalValue{
		Type: "break",
		Exp:  nil,
	}
}

// FOR IN STMT
func (e Evaluator) EvaluateForInStmt(node parser.ForInSatementNode, env *environment.Environment) values.RuntimeValue {

	// Evaluate iterator expression
	iterator := e.EvaluateExpression(node.Iterator, env)

	if IsError(iterator) {
		return iterator
	}

	// Flag for Break statement
	thereIsBreak := false

	if iterator.GetType() == "ArrayValue" {
		// Creating new environment
		env.PushScope()

		iterValues := iterator.(*values.ArrayValue).Value

		for index, value := range iterValues {

			if thereIsBreak == true {
				break
			}

			// Load variables in env on each iteration
			env.DeclareVar(node.LocalVarName, value)

			if node.IndexVarName != "" {
				env.DeclareVar(node.IndexVarName, values.NumberValue{Value: float64(index)})
			}

			// LOOP through for in body!
			for _, stmt := range node.Body {

				result := e.EvaluateStmt(stmt, env)

				if result != nil {
					if result.GetType() == "ErrorValue" || result.GetType() == "return" {
						env.ExitScope()
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
		// Creating new environment
		env.PushScope()

		iterValues := iterator.(values.DictionaryValue).Value

		for index, value := range iterValues {

			if thereIsBreak == true {
				break
			}

			env.DeclareVar(node.LocalVarName, values.StringValue{Value: index, Mutable: false})
			if node.IndexVarName != "" {
				env.DeclareVar(node.IndexVarName, value)
			}

			for _, stmt := range node.Body {

				result := e.EvaluateStmt(stmt, env)

				if result != nil {
					if result.GetType() == "ErrorValue" || result.GetType() == "return" {
						env.ExitScope()
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
		env.ExitScope()
	}

	return values.BooleanValue{Value: true}
}

// STRUCT DECLARATION
func (e Evaluator) EvaluatStructDeclarationStmt(node parser.StructDeclarationNode, env *environment.Environment) values.RuntimeValue {

	rtValue := values.StructValue{}

	rtValue.Properties = node.Properties
	rtValue.Methods = make(map[string]values.RuntimeValue)

	_, err := env.DeclareVar(node.Name, rtValue)

	if err != nil {
		return e.Panic(err.Error(), node.Line, env)
	}

	return values.BooleanValue{Value: true}
}

// FUNCTION DECLARATION
func (e Evaluator) EvaluateFunctionDeclarationStmt(node parser.FunctionDeclarationNode, env *environment.Environment) values.RuntimeValue {

	fn := values.FunctionValue{}

	fn.Body = node.Body
	fn.Parameters = node.Parameters
	fn.Struct = ""
	fn.Environment = env
	fn.Evaluator = e

	_, err := env.DeclareVar(node.Name, fn)

	if err != nil {
		return e.Panic(err.Error(), node.Line, env)
	}

	return values.BooleanValue{Value: true}
}

// IF STMT
func (e Evaluator) EvaluateIfStmt(node parser.IfStatementNode, env *environment.Environment) values.RuntimeValue {

	exp := node.Condition

	evaluatedExp := e.EvaluateExpression(exp, env)

	if IsError(evaluatedExp) {
		return evaluatedExp
	}

	value, err := e.EvaluateImplicitBoolConversion(evaluatedExp)

	if err != nil {
		return e.Panic(err.Error(), node.Line, env)
	}

	if value == true {

		env.PushScope()

		for _, stmt := range node.Body {

			result := e.EvaluateStmt(stmt, env)

			if result != nil && result.GetType() == "ErrorValue" || result.GetType() == "return" || result.GetType() == "break" || result.GetType() == "continue" {
				env.ExitScope()
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

				value, err := e.EvaluateImplicitBoolConversion(evaluatedExp)

				if err != nil {
					return e.Panic(err.Error(), node.Line, env)
				}

				if value == true {

					matched = true

					env.PushScope()

					for _, stmt := range elseif.Body {

						result := e.EvaluateStmt(stmt, env)

						if result != nil && (result.GetType() == "ErrorValue" || result.GetType() == "return" || result.GetType() == "break" || result.GetType() == "continue") {
							env.ExitScope()
							return result
						}
					}
					env.ExitScope()
					return nil
				}
			}
		}

		if matched == false {
			env.PushScope()
			for _, stmt := range node.ElseBody {

				result := e.EvaluateStmt(stmt, env)
				if result != nil && (result.GetType() == "ErrorValue" || result.GetType() == "return" || result.GetType() == "break" || result.GetType() == "continue") {
					env.ExitScope()
					return result
				}
			}
			env.ExitScope()
		}
	}

	return values.BooleanValue{Value: true}
}

// Evaluate a single Expression Statement
// An Expression Statement is statement that is only formed by expressions, like assigments or binary expressions
func (e *Evaluator) EvaluateExpressionStmt(node parser.ExpressionStmtNode, env *environment.Environment) values.RuntimeValue {
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

	_, err := env.DeclareVar(identifier.Value, parsed)

	if err != nil {
		return e.Panic(err.Error(), node.Line, env)
	}

	return values.BooleanValue{Value: true}
}

// Evaluate an expression
func (e *Evaluator) EvaluateExpression(n parser.Exp, env *environment.Environment) values.RuntimeValue {

	if n == nil {
		return nil
	}

	switch n.ExpType() {

	case parser.NodeBinaryExp:
		return e.EvaluateBinaryExpression(n.(parser.BinaryExpNode), env)
	case parser.NodeAnonFunctionDeclaration:
		return e.EvaluateAnonymousFunctionExpression(n.(parser.AnonFunctionDeclarationNode), env)
	case parser.NodeTernaryExp:
		return e.EvaluateTernaryExpression(n.(parser.TernaryExpNode), env)
	case parser.NodeCallExp:
		return e.EvaluateCallExpression(n.(parser.CallExpNode), env)
	case parser.NodeIdentifier:
		node := n.(parser.IdentifierNode)
		lookup, err := env.GetVar(node.Value, node.Line)
		if err != nil {
			return e.Panic(err.Error(), node.Line, env)
		}
		return lookup
	case parser.NodeUnaryExp:
		return e.EvaluateUnaryExpression(n.(parser.UnaryExpNode), env)
	case parser.NodeNumber:
		parsedNumber, _ := strconv.ParseFloat(n.(parser.NumberNode).Value, 2)
		return values.NumberValue{Value: parsedNumber}
	case parser.NodeString:
		return values.StringValue{Value: n.(parser.StringNode).Value}
	case parser.NodeNothing:
		return values.NothingValue{}
	case parser.NodeIndexAccessExp:
		return e.EvaluateIndexAccessExpression(n.(parser.IndexAccessExpNode), env)
	case parser.NodeBoolean:
		return values.BooleanValue{Value: n.(parser.BooleanNode).Value}
	case parser.NodeAssignment:
		return e.EvaluateAssignmentExpression(n.(parser.AssignmentNode), env)
	case parser.NodeSliceExp:
		return e.EvaluateSliceExpression(n.(parser.SliceExpNode), env)
	case parser.NodeMemberExp:
		return e.EvaluateMemberExpression(n.(parser.MemberExpNode), env)
	case parser.NodeObjectInitExp:
		return e.EvaluateObjInitializeExpression(n.(parser.ObjectInitExpNode), env)
	case parser.NodeArrayExp:
		return e.EvaluateArrayExpression(n.(parser.ArrayExpNode), env)
	case parser.NodeDictionaryExp:
		return e.EvaluateDictionaryExpression(n.(parser.DictionaryExpNode), env)
	default:
		return values.ErrorValue{Value: "Unknown expression type"}
	}

}

func (e Evaluator) EvaluateAnonymousFunctionExpression(node parser.AnonFunctionDeclarationNode, env *environment.Environment) values.RuntimeValue {

	fn := values.FunctionValue{}

	fn.Body = node.Body
	fn.Parameters = node.Parameters
	fn.Struct = ""
	fn.Environment = env
	fn.Evaluator = e

	return fn
}

func (e Evaluator) EvaluateTernaryExpression(node parser.TernaryExpNode, env *environment.Environment) values.RuntimeValue {

	condition := e.EvaluateExpression(node.Condition, env)

	if IsError(condition) {
		return condition
	}

	value, err := e.EvaluateImplicitBoolConversion(condition)

	if err != nil {
		return e.Panic(err.Error(), node.Line, env)
	}

	if value {
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
	case *values.ArrayValue:
		if end == nil {
			end = values.NumberValue{Value: float64(len(value.(values.ArrayValue).Value))}
		}

		fn := value.GetProp(&value, "slice")
		ret := fn.(values.NativeFunctionValue).Value([]values.RuntimeValue{init, end})

		return ret
	case values.StringValue:
		if end == nil {
			end = values.NumberValue{Value: float64(len(value.(values.StringValue).Value))}
		}

		fn := value.GetProp(&value, "slice")
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
	structLup, err := env.GetVar(structName, node.Line)

	if err != nil {
		e.Panic(err.Error(), node.Line, env)
	}

	if structLup == nil || structLup.GetType() != "StructValue" {
		e.Panic("Expected struct, got "+structLup.GetType(), node.Line, env)
	}

	// Check if method already exists
	_, exists := structLup.(values.StructValue).Methods[node.Function.Name]

	if exists {
		return e.Panic("Method '"+node.Function.Name+"' already exists in struct '"+structName+"'", node.Line, env)
	}

	// Create function value
	fn := values.FunctionValue{}

	fn.Body = node.Function.Body
	fn.Parameters = node.Function.Parameters
	fn.Struct = structName
	fn.Environment = env

	// Store function in struct
	env.Variables[len(env.Variables)][structName].(values.StructValue).Methods[node.Function.Name] = fn

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

	fn := varValue.GetProp(&varValue, node.Member.Value)

	if fn == nil {
		return e.Panic(varValue.GetType()+" has not member "+node.Member.Value, node.Line, env)
	}

	return fn

}

// Evaluate an object initialization
func (e Evaluator) EvaluateObjInitializeExpression(node parser.ObjectInitExpNode, env *environment.Environment) values.RuntimeValue {

	// Syntax for object initialization is the same as dictionaries
	// So initialize one for obtaining the values
	propDict := node.Value

	// val will be the final object where we will store the values
	val := values.ObjectValue{}

	// Lookup for the base struct
	structLup := e.EvaluateExpression(node.Struct, env)

	if structLup != nil && structLup.GetType() == "ErrorValue" {
		return structLup
	}

	// If when evaluating the struct it is not a struct, return error
	if _, ok := structLup.(values.StructValue); !ok {
		return e.Panic("You only can initialize objects of structs, not of "+structLup.GetType(), node.Line, env)
	}

	val.Struct = structLup.(values.StructValue)
	val.Value = make(map[string]values.RuntimeValue)

	// Create a map with the properties defined in the struct
	// So we can check if the property exists with less complexity
	structProperties := make(map[string]bool)

	// Initialize each property of the object with nothing
	// And at the same type fill the map of properties created earlier
	for _, prop := range val.Struct.Properties {
		val.Value[prop] = values.NothingValue{}
		structProperties[prop] = true
	}

	// Evaluate expressions and set values
	// For each property defined in the initialization, which is parsed as a dictionary
	// Check if that property exists in the struct by using the map created earlier
	for key, exp := range propDict.Value {

		if _, ok := structProperties[key]; !ok {
			return e.Panic("Unknown property "+key, node.Line, env)
		}

		value := e.EvaluateExpression(exp, env) // Evaluate value

		if value.GetType() == "ErrorValue" {
			return value
		}

		// Add the value to the map of properties of the object
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
	case *values.ArrayValue:
		iToInt, _ := strconv.Atoi(i)
		if iToInt < 0 {
			iToInt = len(val.Value) + iToInt
		}
		if iToInt >= len(val.Value) {
			return e.Panic("Index "+i+" out of range", node.Line, env)
		}
		return val.Value[iToInt]
	case values.StringValue:
		iToInt, _ := strconv.Atoi(i)
		return values.StringValue{Value: string(val.Value[iToInt]), Mutable: false}
	case values.DictionaryValue:
		item, exists := val.Value[i]

		if !exists {
			return e.Panic("Undefined key '"+i, node.Line, env)
		}

		return item

	default:
		return e.Panic("Only arrays and dictionaries can be accessed by index", node.Line, env)
	}

}

func (e Evaluator) EvaluateArrayExpression(node parser.ArrayExpNode, env *environment.Environment) values.RuntimeValue {

	rtvalue := &values.ArrayValue{}
	rtvalue.Value = make([]values.RuntimeValue, 0)

	for _, exp := range node.Value {
		rtvalue.Value = append(rtvalue.Value, e.EvaluateExpression(exp, env))
	}

	return rtvalue

}

func (e *Evaluator) EvaluateCallExpression(node parser.CallExpNode, env *environment.Environment) values.RuntimeValue {

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

		val := fn.Value(evaluatedArgs)

		// NATIVE FUNCTION RETURNS ERROR VALUES WITH DIFERENT FORMAT SO CREATE A NEW PANIC
		// TODO: make native functions return errors like environment functions with return values = (RuntimeValue, error)
		if val.GetType() == "ErrorValue" {
			return e.Panic(val.GetStr(), node.Line, env)
		}

		return val
	case values.FunctionValue:

		// Create new environment
		env.PushScope()

		// Set this
		if fn.Struct != "" {
			env.DeclareVar("this", fn.StructObjRef)
		}

		for index, arg := range evaluatedArgs {
			env.ForceDeclare(fn.Parameters[index], arg)
		}

		var result values.RuntimeValue

		for _, stmt := range fn.Body {
			result = e.EvaluateStmt(stmt, env)

			if result != nil && result.GetType() == "ErrorValue" {
				env.ExitScope()
				return result
			}

			signal, isSignal := result.(values.SignalValue)

			if isSignal && signal.Type == "return" {

				if signal.Value.GetType() == "ErrorValue" {
					env.ExitScope()
					return signal
				}

				env.ExitScope()

				return signal.Value
			}

		}
		env.ExitScope()
		return values.NothingValue{}

	default:
		return e.Panic("Only functions can be called not "+calle.GetType(), node.Line, env)
	}

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
	if node.Left.ExpType() == parser.NodeIndexAccessExp {
		indexes = append(e.ResolveIndexAccessChain(node.Left.(parser.IndexAccessExpNode), env), indexes...)
	} else if node.Left.ExpType() == parser.NodeIdentifier {
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

	if left.ExpType() == parser.NodeIndexAccessExp {
		chain := e.ResolveIndexAccessChain(left.(parser.IndexAccessExpNode), env)
		_, err := env.ModifyIndexValue(left.(parser.IndexAccessExpNode), right, chain)

		if err != nil {
			return e.Panic(err.Error(), node.Line, env)
		}

		return right
	} else {
		_, err := env.SetVar(left, right)

		if err != nil {
			return e.Panic(err.Error(), node.Line, env)
		}
	}

	return right
}

func (e Evaluator) EvaluateBinaryExpression(node parser.BinaryExpNode, env *environment.Environment) values.RuntimeValue {

	left := e.EvaluateExpression(node.Left, env)
	right := e.EvaluateExpression(node.Right, env)

	if left.GetType() == "ErrorValue" {
		return left
	}
	if right.GetType() == "ErrorValue" {
		return right
	}

	op := node.Operator

	type1 := left.GetType()
	type2 := right.GetType()

	equalTypes := type1 == type2

	if op == "+" {
		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "StringValue" {
			return values.StringValue{Value: left.GetStr() + right.GetStr(), Mutable: false}
		} else if type1 == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() + right.GetNumber()}
		} else {
			return e.Panic("Cant use operator "+op+" with type "+type1, node.Line, env)
		}

	} else if op == "-" {
		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() - right.GetNumber()}
		} else {
			return e.Panic("Cant use operator "+op+" with type "+type1, node.Line, env)
		}
	} else if op == "*" {
		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() * right.GetNumber()}
		} else {
			return e.Panic("Cant use operator "+op+" with type "+type1, node.Line, env)
		}
	} else if op == "/" {
		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "NumberValue" {
			return values.NumberValue{Value: left.GetNumber() / right.GetNumber()}
		} else {
			return e.Panic("Cant use operator "+op+" with type "+type1, node.Line, env)
		}
	} else if op == "==" {

		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "StringValue" {
			return values.BooleanValue{Value: left.GetStr() == right.GetStr()}
		} else if type1 == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() == right.GetNumber()}
		} else if type1 == "BooleanValue" {
			return values.BooleanValue{Value: left.GetBool() == right.GetBool()}
		} else {
			return e.Panic("Cant use operator "+op+" with type "+type1, node.Line, env)
		}
	} else if op == ">" {

		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() > right.GetNumber()}
		} else {
			return e.Panic("Operator "+op+" only can be used with numbers, not with type "+type1, node.Line, env)
		}
	} else if op == "<" {

		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() < right.GetNumber()}
		} else {
			return e.Panic("Operator "+op+" only can be used with numbers, not with type "+type1, node.Line, env)
		}
	} else if op == "<=" {

		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() <= right.GetNumber()}
		} else {
			return e.Panic("Operator "+op+" only can be used with numbers, not with type "+type1, node.Line, env)
		}
	} else if op == ">=" {

		if !equalTypes {
			return e.Panic("Type mismatch: "+type1+" and "+type2, node.Line, env)
		}

		if type1 == "NumberValue" {
			return values.BooleanValue{Value: left.GetNumber() >= right.GetNumber()}
		} else {
			return e.Panic("Operator "+op+" only can be used with numbers, not with type "+type1, node.Line, env)
		}
	} else if op == "and" {

		leftValue, err := e.EvaluateImplicitBoolConversion(left)

		if err != nil {
			return e.Panic(err.Error(), node.Line, env)
		}

		rightValue, err := e.EvaluateImplicitBoolConversion(right)

		if err != nil {
			return e.Panic(err.Error(), node.Line, env)
		}

		return values.BooleanValue{Value: leftValue && rightValue}

	} else if op == "or" {

		leftValue, err := e.EvaluateImplicitBoolConversion(left)

		if err != nil {
			return e.Panic(err.Error(), node.Line, env)
		}

		rightValue, err := e.EvaluateImplicitBoolConversion(right)

		if err != nil {
			return e.Panic(err.Error(), node.Line, env)
		}

		return values.BooleanValue{Value: leftValue || rightValue}

	}

	return nil

}

func (e Evaluator) EvaluateUnaryExpression(node parser.UnaryExpNode, env *environment.Environment) values.RuntimeValue {

	exp := e.EvaluateExpression(node.Right, env)

	if node.Operator == "-" && exp.GetType() == "NumberValue" {
		return values.NumberValue{Value: -exp.GetNumber()}
	} else if node.Operator == "not" {
		return values.BooleanValue{Value: !exp.GetBool()}
	} else {
		return nil
	}
}

func (e Evaluator) EvaluateImplicitBoolConversion(value values.RuntimeValue) (bool, error) {
	switch val := value.(type) {
	case values.NumberValue:
		return val.GetBool(), nil
	case values.StringValue:
		return val.GetBool(), nil
	case values.BooleanValue:
		return val.GetBool(), nil
	case values.ArrayValue:
		return val.GetBool(), nil
	default:
		return false, errors.New("Cannot convert " + value.GetType() + " to boolean")
	}
}

func (e Evaluator) ExecuteCallback(fn interface{}, env interface{}, args []interface{}) interface{} {

	fnValue := fn.(values.FunctionValue)

	// parentEnvironment := env.(*environment.Environment)

	childEnvironment := environment.NewEnvironment()

	for i, paramName := range fnValue.Parameters {
		if i+1 < len(args) {
			break
		}

		switch val := args[i].(type) {
		case values.ObjectValue:
			childEnvironment.Variables[len(childEnvironment.Variables)-1][paramName] = values.RuntimeValue(val)
		case values.NumberValue:
			childEnvironment.Variables[len(childEnvironment.Variables)-1][paramName] = values.RuntimeValue(val)
		case values.StringValue:
			childEnvironment.Variables[len(childEnvironment.Variables)-1][paramName] = values.RuntimeValue(val)
		case values.BooleanValue:
			childEnvironment.Variables[len(childEnvironment.Variables)-1][paramName] = values.RuntimeValue(val)
		case values.ArrayValue:
			childEnvironment.Variables[len(childEnvironment.Variables)-1][paramName] = values.RuntimeValue(val)
		}
	}

	var result values.RuntimeValue

	for _, stmt := range fnValue.Body {
		result = e.EvaluateStmt(stmt, childEnvironment)

		if result != nil && result.GetType() == "ErrorValue" {
			return result
		}

		signal, isSignal := result.(values.SignalValue)

		if isSignal && signal.Type == "return" {

			return signal.Value
		}

	}
	return nil
}
