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
func (e Evaluator) Evaluate(env *environment.Environment) *environment.Environment {

	for _, node := range e.Nodes {

		ret := e.EvaluateStmt(node, env)

		// If the return value is an ErrorValue
		if ret.GetType() == values.ErrorType {

			e.PrintError(ret.(values.ErrorValue))
			os.Exit(1)
		}

	}
	return env
}

// Evaluate a single Statement node
func (e Evaluator) EvaluateStmt(n parser.Stmt, env *environment.Environment) values.RuntimeValue {
	switch n.StmtType() {
	case parser.NodeExpStmt:
		return e.EvaluateExpression(n.(parser.ExpressionStmtNode).Expression, env)
	case parser.NodeVarDeclaration:
		return e.EvaluateVarDeclaration(n.(parser.VarDeclarationNode), env)
	case parser.NodeIfStatement:
		//init := timer.Init()

		//timer.add("if_stmt", init)
		return e.EvaluateIfStmt(n.(parser.IfStatementNode), env)
	case parser.NodeForInStatement:
		return e.EvaluateForInStmt(n.(parser.ForInSatementNode), env)
	case parser.NodeFunctionDeclaration:
		return e.EvaluateFunctionDeclarationStmt(n.(parser.FunctionDeclarationNode), env)
	case parser.NodeReturnStatement:
		return e.EvaluateReturnNode(n.(parser.ReturnNode), env)
	case parser.NodeLoopStatement:
		return e.EvaluateLoopStmt(n.(parser.LoopStmtNode), env)
	case parser.NodeStructMethodDeclaration:
		return e.EvaluateStructMethodExpression(n.(parser.StructMethodDeclarationNode), env)
	case parser.NodeBreakStatement:
		return e.EvaluateBreakNode()
	case parser.NodeContinueStatement:
		return e.EvaluateContinueNode(n.(parser.ContinueNode), env)
	case parser.NodeTryCatchStatement:
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

		if _, ok := env.ImportChain[node.Path]; ok {
			return e.Panic(values.CircularImportError, "Circular import with module: "+node.Path, line, env)
		}

		// Do not add .ev to modules ;)
		path := node.Path + ".ev"

		// Read file, parse and evaluate
		var source string = utils.ReadFile(path)

		tokens := lexer.Tokenize(source)
		ast := parser.NewParser(tokens).GetAST()

		// Create new environment for the module with the parent environment
		envForModule := environment.NewEnvironment()
		native.SetupEnvironment(envForModule)
		envForModule.ImportChain = env.ImportChain
		envForModule.ImportChain[env.ModuleName] = true
		envForModule.ModuleName = node.Path

		eval := Evaluator{Nodes: ast}
		importEnv := eval.Evaluate(envForModule)
		// Get the created environment after evaluate the module
		// Get all the variables loaded and load into the actual environment
		// using a namespace
		env.DeclareVar(node.Alias, values.NamespaceValue{Value: importEnv.Variables})

	}

	return values.BoolValue{Value: true}
}

// LOOP STATEMENT
func (e Evaluator) EvaluateLoopStmt(node parser.LoopStmtNode, env *environment.Environment) values.RuntimeValue {

	// Creating new environment

	for {

		loopenv := environment.NewScopeEnv(env, 0)

		// Loop through body
		for _, stmt := range node.Body {

			ret := e.EvaluateStmt(stmt, loopenv)

			t := ret.GetType()

			if t == values.ErrorType || t == values.ReturnType {
				return ret
			} else if t == values.BreakType {
				return values.NothingValue{}
			} else if t == values.ContinueType {
				break
			}
		}
	}
}

// TRY CATCH
func (e Evaluator) EvaluateTryCatchNode(node parser.TryCatchNode, env *environment.Environment) values.RuntimeValue {

	body := node.Body
	catch := node.Catch

	// Loop through body
	for _, stmt := range body {

		ret := e.EvaluateStmt(stmt, env)

		if ret.GetType() == values.ReturnType {
			if node.Finally != nil {
				e.EvaluateFinallyBlock(node.Finally, env)
			}
			return ret
		}

		if ret.GetType() == values.ErrorType {

			env.DeclareVar("error", ret.(values.ErrorValue).Object)

			// CATCH BODY
			for _, cstmt := range catch {

				result := e.EvaluateStmt(cstmt, env)

				// Error inside the catch lol
				if result.GetType() == values.ErrorType || result.GetType() == values.ReturnType {
					if node.Finally != nil {
						e.EvaluateFinallyBlock(node.Finally, env)
					}

					return result
				} else if result.GetType() == values.BreakType {

					return values.BoolValue{Value: true}
				} else if result.GetType() == values.ContinueType {
					continue
				}
			}

			if node.Finally != nil {
				e.EvaluateFinallyBlock(node.Finally, env)
			}

			return values.BoolValue{Value: true}
		}
	}
	if node.Finally != nil {
		e.EvaluateFinallyBlock(node.Finally, env)
	}

	return values.BoolValue{Value: true}
}

func (e Evaluator) EvaluateFinallyBlock(stmt []parser.Stmt, env *environment.Environment) values.RuntimeValue {

	scope := environment.NewScopeEnv(env, 0)

	for _, stmt := range stmt {
		res := e.EvaluateStmt(stmt, scope)

		if res.GetType() == values.ErrorType {
			return res
		}
	}

	return values.NothingValue{}

}

// RETURN STMT
func (e Evaluator) EvaluateReturnNode(node parser.ReturnNode, env *environment.Environment) values.RuntimeValue {

	eval := e.EvaluateExpression(node.Right, env)
	if eval.GetType() == values.ErrorType {
		return eval
	}

	return values.ReturnValue{
		Value: eval,
	}
}

// CONTINUE
func (e Evaluator) EvaluateContinueNode(node parser.ContinueNode, env *environment.Environment) values.RuntimeValue {
	return values.ContinueValue{}
}

// BREAk
func (e Evaluator) EvaluateBreakNode() values.RuntimeValue {
	return values.BreakValue{}
}

// FOR IN STMT
func (e Evaluator) EvaluateForInStmt(node parser.ForInSatementNode, env *environment.Environment) values.RuntimeValue {

	// // Evaluate iterator expression
	iterator := e.EvaluateExpression(node.Iterator, env)

	if iterator.GetType() == values.ErrorType {
		return iterator
	}

	// Flag for Break statement
	thereIsBreak := false

	if iterator.GetType() == values.ArrayType {
		// Creating new environment

		iterValues := iterator.(*values.ArrayValue).Value

		for index, value := range iterValues {

			loopenv := environment.NewScopeEnv(env, 2)

			if thereIsBreak == true {
				break
			}

			// Load variables in env on each iteration
			loopenv.ForceDeclare(node.LocalVarName, value)

			if node.IndexVarName != "" {
				loopenv.ForceDeclare(node.IndexVarName, values.NumberValue{Value: float64(index)})
			}

			// LOOP through for in body!
			for _, stmt := range node.Body {

				result := e.EvaluateStmt(stmt, loopenv)

				if result.GetType() == values.ErrorType || result.GetType() == values.ReturnType {

					return result
				} else if result.GetType() == values.BreakType {
					thereIsBreak = true
					break
				} else if result.GetType() == values.ContinueType {
					break
				}

			}
		}

	} else if iterator.GetType() == values.DictionaryType {
		// Creating new environment

		loopenv := environment.NewScopeEnv(env, 2)

		iterValues := iterator.(*values.DictionaryValue).Value

		for index, value := range iterValues {

			if thereIsBreak == true {
				break
			}

			loopenv.ForceDeclare(node.LocalVarName, values.StringValue{Value: index})
			if node.IndexVarName != "" {
				loopenv.ForceDeclare(node.IndexVarName, value)
			}

			for _, stmt := range node.Body {

				result := e.EvaluateStmt(stmt, loopenv)

				if result.GetType() == values.ErrorType || result.GetType() == values.ReturnType {

					return result
				} else if result.GetType() == values.BreakType {
					thereIsBreak = true
					break
				} else if result.GetType() == values.ContinueType {
					break
				}

			}

		}

	}

	return values.BoolValue{Value: true}
}

// STRUCT DECLARATION
func (e Evaluator) EvaluatStructDeclarationStmt(node parser.StructDeclarationNode, env *environment.Environment) values.RuntimeValue {

	rtValue := values.StructValue{}

	rtValue.Properties = node.Properties
	rtValue.Methods = make(map[string]values.RuntimeValue)
	rtValue.Name = node.Name

	err := env.DeclareVar(node.Name, rtValue)

	if err != nil {
		return e.Panic(values.IdentifierError, err.Error(), node.Line, env)
	}

	return values.BoolValue{Value: true}
}

// FUNCTION DECLARATION
func (e Evaluator) EvaluateFunctionDeclarationStmt(node parser.FunctionDeclarationNode, env *environment.Environment) values.RuntimeValue {

	fn := values.FunctionValue{}

	fn.Body = node.Body
	fn.Parameters = node.Parameters
	fn.Struct = ""
	fn.Evaluator = e
	fnenv := env
	fn.Environment = fnenv

	err := env.DeclareVar(node.Name, fn)

	if err != nil {
		return e.Panic(values.IdentifierError, err.Error(), node.Line, env)
	}

	return values.BoolValue{Value: true}
}

// IF STMT
func (e Evaluator) EvaluateIfStmt(node parser.IfStatementNode, env *environment.Environment) values.RuntimeValue {

	// litter.Dump(node)
	// os.Exit(1)

	evaluatedExp := e.EvaluateExpression(node.Condition, env)

	if evaluatedExp.GetType() == values.ErrorType {
		return evaluatedExp
	}

	value, err := e.EvaluateImplicitBoolConversion(evaluatedExp)

	if err != nil {
		return e.Panic(values.InvalidConversionError, err.Error(), node.Line, env)
	}

	if value == true {
		ifenv := environment.NewScopeEnv(env, 0)

		for _, stmt := range node.Body {

			result := e.EvaluateStmt(stmt, ifenv)

			t := result.GetType()

			if t == values.ErrorType || t == values.ReturnType || t == values.BreakType || t == values.ContinueType {

				return result
			}

		}

	} else {

		matched := false

		if len(node.ElseIf) > 0 {
			for _, elseif := range node.ElseIf {

				if matched == true {
					break
				}

				exp := e.EvaluateExpression(elseif.Condition, env)

				if exp.GetType() == values.ErrorType {
					return exp
				}

				value, err := e.EvaluateImplicitBoolConversion(exp)

				if err != nil {
					return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
				}

				if value == true {

					matched = true

					ifenv := environment.NewScopeEnv(env, 0)

					for _, stmt := range elseif.Body {

						result := e.EvaluateStmt(stmt, ifenv)

						if result.GetType() == values.ErrorType || result.GetType() == values.ReturnType || result.GetType() == values.BreakType || result.GetType() == values.ContinueType {

							return result
						}
					}

					return values.BoolValue{Value: true}
				}
			}
		}

		// if node.ElseBody != nil && matched == false {
		if node.ElseBody != nil && matched == false {

			ifenv := environment.NewScopeEnv(env, 0)

			for _, stmt := range node.ElseBody {

				result := e.EvaluateStmt(stmt, ifenv)
				if result.GetType() == values.ErrorType || result.GetType() == values.ReturnType || result.GetType() == values.BreakType || result.GetType() == values.ContinueType {

					return result
				}
			}

		}
	}

	return values.NothingValue{}
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

	parsedtype := parsed.GetType()

	if parsedtype == values.ErrorType {
		return parsed
	}

	err := env.DeclareVar(identifier.Value, parsed)

	if err != nil {
		return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
	}

	return values.NothingValue{}
}

// Evaluate an expression
func (e Evaluator) EvaluateExpression(n parser.Exp, env *environment.Environment) values.RuntimeValue {

	switch n.ExpType() {

	case parser.NodeBinaryExp:
		// init := timer.Init()
		// v := e.EvaluateBinaryExpression(n.(parser.BinaryExpNode), env)
		// timer.Add("bin_exp", init)
		return e.EvaluateBinaryExpression(n.(parser.BinaryExpNode), env)
	case parser.NodeBinaryComparisonExp:
		// init := timer.Init()
		v := e.EvaluateBinaryComparisonExpression(n.(parser.BinaryComparisonExpNode), env)
		// timer.Add("bin_exp", init)
		return v
	case parser.NodeBinaryLogicExp:
		// init := timer.Init()
		v := e.EvaluateBinaryLogicExpression(n.(parser.BinaryLogicExpNode), env)
		// timer.Add("bin_exp", init)
		return v
	case parser.NodeTernaryExp:
		return e.EvaluateTernaryExpression(n.(parser.TernaryExpNode), env)
	case parser.NodeNumber:
		return values.NumberValue{Value: n.(parser.NumberNode).Value}
	case parser.NodeCallExp:
		return e.EvaluateCallExpression(n.(parser.CallExpNode), env)
	case parser.NodeAnonFunctionDeclaration:
		return e.EvaluateAnonymousFunctionExpression(n.(parser.AnonFunctionDeclarationNode), env)
	case parser.NodeIdentifier:
		node := n.(parser.IdentifierNode)
		lookup, err := env.GetVar(node.Value)
		if err != nil {
			return e.Panic(values.IdentifierError, err.Error(), node.Line, env)
		}
		return lookup
	case parser.NodeUnaryExp:
		return e.EvaluateUnaryExpression(n.(parser.UnaryExpNode), env)
	case parser.NodeString:
		return values.StringValue{Value: n.(parser.StringNode).Value}
	case parser.NodeNothing:
		return values.NothingValue{}
	case parser.NodeIndexAccessExp:
		return e.EvaluateIndexAccessExpression(n.(parser.IndexAccessExpNode), env)
	case parser.NodeBoolean:
		return values.BoolValue{Value: n.(parser.BooleanNode).Value}
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
		// litter.Dump(n)
		return e.Panic(values.RuntimeError, "Unknown Expression Type", 0, env)
	}

}

func (e Evaluator) EvaluateAnonymousFunctionExpression(node parser.AnonFunctionDeclarationNode, env *environment.Environment) values.RuntimeValue {

	fn := values.FunctionValue{}

	fn.Body = node.Body
	fn.Parameters = node.Parameters
	fn.Struct = ""

	// fnenv := environment.NewEnvironment()
	// native.SetupEnvironment(fnenv)

	fn.Environment = env
	fn.Evaluator = e

	return fn
}

func (e Evaluator) EvaluateTernaryExpression(node parser.TernaryExpNode, env *environment.Environment) values.RuntimeValue {

	condition := e.EvaluateExpression(node.Condition, env)

	if condition.GetType() == values.ErrorType {
		return condition
	}

	value, err := e.EvaluateImplicitBoolConversion(condition)

	if err != nil {
		return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
	}

	if value {
		return e.EvaluateExpression(node.Left, env)
	} else {
		return e.EvaluateExpression(node.Right, env)
	}
}

func (e Evaluator) EvaluateSliceExpression(node parser.SliceExpNode, env *environment.Environment) values.RuntimeValue {

	value := e.EvaluateExpression(node.Left, env)

	var init values.RuntimeValue = nil
	var end values.RuntimeValue = nil

	if node.From != nil {
		init = e.EvaluateExpression(node.From, env)
		if init.GetType() == values.ErrorType || value.GetType() == values.ErrorType {
			return init
		}
	}

	if node.To != nil {
		end = e.EvaluateExpression(node.To, env)
		if end.GetType() == values.ErrorType {
			return end
		}
	}

	switch value.GetType() {
	case values.ArrayType:

		fn, err := value.(*values.ArrayValue).GetProp("slice")

		if err != nil {
			return e.Panic(values.PropertyError, err.Error(), node.Line, env)
		}

		var ret values.RuntimeValue

		if node.To == nil {
			ret = fn.(values.NativeFunctionValue).Value([]values.RuntimeValue{init})
		} else {
			ret = fn.(values.NativeFunctionValue).Value([]values.RuntimeValue{init, end})
		}

		if ret.GetType() == values.ErrorType {
			return e.Panic(ret.(values.ErrorValue).ErrorType, ret.(values.ErrorValue).Value, node.Line, env)
		}

		return ret
	case values.StringType:

		fn, err := value.GetProp("slice")

		if err != nil {
			return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
		}
		var ret values.RuntimeValue

		if node.To == nil {
			ret = fn.(values.NativeFunctionValue).Value([]values.RuntimeValue{init})
		} else {
			ret = fn.(values.NativeFunctionValue).Value([]values.RuntimeValue{init, end})
		}

		if ret.GetType() == values.ErrorType {
			return e.Panic(ret.(values.ErrorValue).ErrorType, ret.(values.ErrorValue).Value, node.Line, env)
		}
		return ret
	default:
		return values.ErrorValue{Value: "Expected array or string"}
	}

}

// Declaration of a struct method
func (e Evaluator) EvaluateStructMethodExpression(node parser.StructMethodDeclarationNode, env *environment.Environment) values.RuntimeValue {

	structName := node.Struct

	// Check if struct exists and is a struct
	structLup, err := env.GetVar(structName)

	if err != nil {
		e.Panic(values.IdentifierError, err.Error(), node.Line, env)
	}

	if structLup.GetType() != values.StructType {
		e.Panic(values.TypeError, "Expected struct, got "+structLup.GetType().String(), node.Line, env)
	}

	// Check if method already exists
	_, exists := structLup.(values.StructValue).Methods[node.Function.Name]

	if exists {
		return e.Panic(values.RuntimeError, "Method '"+node.Function.Name+"' already exists in struct '"+structName+"'", node.Line, env)
	}

	// Create function value
	fn := values.FunctionValue{}

	fn.Body = node.Function.Body
	fn.Parameters = node.Function.Parameters
	fn.Struct = structName
	fn.Environment = env

	// Store function in struct
	env.Variables[structName].(values.StructValue).Methods[node.Function.Name] = fn

	return values.NothingValue{}
}

// Evaluate a member expression
func (e Evaluator) EvaluateMemberExpression(node parser.MemberExpNode, env *environment.Environment) values.RuntimeValue {

	varValue := e.EvaluateExpression(node.Left, env)

	if varValue.GetType() == values.ErrorType {
		return varValue
	}
	//init := timer.Init()
	fn, err := varValue.GetProp(node.Member)
	//timer.add("get_prop", init)

	if err != nil {
		return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
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

	if structLup.GetType() == values.ErrorType {
		return structLup
	}

	// If when evaluating the struct it is not a struct, return error
	if _, ok := structLup.(values.StructValue); !ok {
		return e.Panic(values.RuntimeError, "You only can initialize objects of structs, not of "+structLup.GetType().String(), node.Line, env)
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
			return e.Panic(values.RuntimeError, "Unknown property "+key, node.Line, env)
		}

		value := e.EvaluateExpression(exp, env) // Evaluate value

		if value.GetType() == values.ErrorType {
			return value
		}

		// Add the value to the map of properties of the object
		val.Value[key] = value
	}

	return &val
}

func (e Evaluator) EvaluateDictionaryExpression(node parser.DictionaryExpNode, env *environment.Environment) values.RuntimeValue {

	dict := values.DictionaryValue{}
	dictValue := make(map[string]values.RuntimeValue)

	for key, exp := range node.Value {
		value := e.EvaluateExpression(exp, env)
		if value.GetType() == values.ErrorType {
			return value
		}
		dictValue[key] = value
	}

	dict.Value = dictValue

	return &dict
}

// Evaluate an index access
func (e Evaluator) EvaluateIndexAccessExpression(node parser.IndexAccessExpNode, env *environment.Environment) values.RuntimeValue {

	// // Obtenemos el valor final del valor base
	identifier := e.EvaluateExpression(node.Left, env)

	if identifier.GetType() == values.ErrorType {
		return identifier
	}

	// Obtenemos el valor final del indice
	index := e.EvaluateExpression(node.Index, env)

	if index.GetType() == values.ErrorType {
		return identifier
	}

	// El indice puede ser numeric si es un array, o un string si se trata de un diccionario
	// En ambos casos lo trataremos como un string y si es necesario se convertira a int

	var i string = index.GetString()

	switch identifier.GetType() {
	case values.ArrayType:
		val := identifier.(*values.ArrayValue)
		iToInt, _ := strconv.Atoi(i)
		if iToInt < 0 {
			iToInt = len(val.Value) + iToInt
		}
		if iToInt >= len(val.Value) {
			return e.Panic(values.InvalidIndexError, "Index "+i+" out of range", node.Line, env)
		}
		return val.Value[iToInt]
	case values.StringType:
		val := identifier.(values.StringValue).Value
		iToInt, _ := strconv.Atoi(i)
		if iToInt >= len(val) {
			return e.Panic(values.InvalidIndexError, "Index "+i+" out of range", node.Line, env)
		}
		return values.StringValue{Value: string(val[iToInt])}
	case values.DictionaryType:
		val := identifier.(*values.DictionaryValue)
		item, exists := val.Value[i]

		if !exists {
			return e.Panic(values.RuntimeError, "Undefined key '"+i, node.Line, env)
		}

		return item

	default:
		return e.Panic(values.RuntimeError, "Only arrays and dictionaries can be accessed by index", node.Line, env)
	}

}

func (e Evaluator) EvaluateArrayExpression(node parser.ArrayExpNode, env *environment.Environment) values.RuntimeValue {

	rtvalue := values.ArrayValue{}
	rtvalue.Value = make([]values.RuntimeValue, 0)

	for _, exp := range node.Value {
		rtvalue.Value = append(rtvalue.Value, e.EvaluateExpression(exp, env))
	}

	return &rtvalue

}

func (e *Evaluator) EvaluateCallExpression(node parser.CallExpNode, env *environment.Environment) values.RuntimeValue {

	//init := timer.Init()
	evaluatedArgs := make([]values.RuntimeValue, 0, len(node.Args))

	for _, arg := range node.Args {
		value := e.EvaluateExpression(arg, env)
		if value.GetType() == values.ErrorType {
			return value
		}
		evaluatedArgs = append(evaluatedArgs, value)

	}

	//timer.add("evaluate_args", init)

	//init = timer.Init()
	calle := e.EvaluateExpression(node.Name, env)
	//timer.add("eval_calee", init)

	if calle.GetType() == values.ErrorType {
		return calle
	}

	switch calle.GetType() {

	case values.NativeFunctionType:
		//init := timer.Init()

		val := calle.(values.NativeFunctionValue).Value(evaluatedArgs)

		if val.GetType() == values.ErrorType {
			return e.Panic(val.(values.ErrorValue).ErrorType, val.GetString(), node.Line, env)
		}

		//timer.add("native_fn", init)
		return val
	case values.FunctionType:
		//init := timer.Init()
		fn := calle.(values.FunctionValue)

		fnEnv := environment.NewScopeEnv(fn.Environment.(*environment.Environment), len(node.Args))
		// Create new environment

		e.AddToCallStack(node.Line, env)

		for index, arg := range evaluatedArgs {
			fnEnv.ForceDeclare(fn.Parameters[index], arg)
		}
		// Set this
		if fn.Struct != "" {
			fnEnv.ForceDeclare("this", fn.StructObjRef)
		}
		var result values.RuntimeValue
		// timer.Add("fn_pre_exec", init)

		for _, stmt := range fn.Body {
			result = e.EvaluateStmt(stmt, fnEnv)

			if result.GetType() == values.ErrorType {

				e.RemoveToCallStack()
				return result
			}

			if result.GetType() == values.ReturnType {

				e.RemoveToCallStack()
				return result.(values.ReturnValue).Value
			}

		}

		e.RemoveToCallStack()
		return values.NothingValue{}

	default:
		return e.Panic(values.RuntimeError, "Only functions can be called not "+calle.GetType().String(), node.Line, env)
	}

}

// Evaluate an Assignment Expression
func (e Evaluator) EvaluateAssignmentExpression(node parser.AssignmentNode, env *environment.Environment) values.RuntimeValue {
	// Left can be an expression, like a index access, or member, etc.
	// We will make the env function to resolve this.
	left := node.Left

	// Evaluate the expression of the right side
	right := e.EvaluateExpression(node.Right, env)

	if right.GetType() == values.ErrorType {
		return right
	}

	if left.ExpType() == parser.NodeIndexAccessExp {
		expNode := left.(parser.IndexAccessExpNode)
		val := e.EvaluateExpression(expNode.Left, env)

		if val.GetType() == values.ErrorType {
			return val
		}

		if val.GetType() == values.ArrayType {
			index := e.EvaluateExpression(expNode.Index, env)

			if index.GetType() == values.ErrorType {
				return index
			}

			if index.GetType() != values.NumberType {
				return e.Panic(values.RuntimeError, "Invalid array index", node.Line, env)
			}

			finalIndex := int(index.GetNumber())

			if finalIndex < 0 {
				finalIndex = len(val.(*values.ArrayValue).Value) + finalIndex
			}

			if finalIndex >= len(val.(*values.ArrayValue).Value) {
				return e.Panic(values.RuntimeError, "Invalid array index or out of bounds with index: "+fmt.Sprint(finalIndex), node.Line, env) // fmt.Sprintf("Invalid array index or out of bounds with index: %d", node.Line, env)
			}

			val.(*values.ArrayValue).Value[int(index.GetNumber())] = right
		} else if val.GetType() == values.DictionaryType {
			key := e.EvaluateExpression(expNode.Index, env)

			if key.GetType() == values.ErrorType {
				return key
			}

			if key.GetType() != values.StringType {
				return e.Panic(values.RuntimeError, "Invalid dictionary key", node.Line, env)
			}

			val.(*values.DictionaryValue).Value[key.GetString()] = right
		}

	} else if left.ExpType() == parser.NodeMemberExp {
		expNode := left.(parser.MemberExpNode)
		val := e.EvaluateExpression(expNode.Left, env)

		if val.GetType() == values.ErrorType {
			return val
		}

		if val.GetType() != values.ObjectType {
			return e.Panic(values.RuntimeError, "Invalid object assignment", node.Line, env)
		}

		val.(*values.ObjectValue).Value[expNode.Member] = right

	} else if left.ExpType() == parser.NodeIdentifier {
		err := env.SetVar(left.(parser.IdentifierNode).Value, right)

		if err != nil {
			return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
		}
	} else {
		return e.Panic(values.RuntimeError, "Invalid assignment", node.Line, env)
	}

	return right
}

func (e Evaluator) EvaluateBinaryExpression(node parser.BinaryExpNode, env *environment.Environment) values.RuntimeValue {

	left := e.EvaluateExpression(node.Left, env)
	right := e.EvaluateExpression(node.Right, env)

	type1 := left.GetType()
	type2 := right.GetType()

	if type1 == values.ErrorType {
		return left
	}
	if type2 == values.ErrorType {
		return right
	}

	equalTypes := type1 == type2

	if !equalTypes {
		return e.Panic(values.RuntimeError, "Type mismatch: "+type1.String()+" and "+type2.String(), node.Line, env)
	}

	if node.Operator == parser.OperatorAdd {

		if type1 == values.NumberType {
			// val := left.(values.NumberValue)
			// val.Value = val.Value + right.(values.NumberValue).Value
			// return val
			return values.NumberValue{Value: left.GetNumber() + right.GetNumber()}
		} else if type1 == values.StringType {
			return values.StringValue{Value: left.(values.StringValue).Value + right.(values.StringValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Cant use operator + with type "+type1.String(), node.Line, env)
		}

	} else if node.Operator == parser.OperatorSubtract {

		if type1 == values.NumberType {
			val := left.(values.NumberValue)
			val.Value = val.Value - right.(values.NumberValue).Value
			return val
			// return values.NumberValue{Value: left.GetNumber() - right.GetNumber()}
		} else {
			return e.Panic(values.RuntimeError, "Cant use operator - with type "+type1.String(), node.Line, env)
		}
	} else if node.Operator == parser.OperatorMultiply {

		if type1 == values.NumberType {
			return values.NumberValue{Value: left.(values.NumberValue).Value * right.(values.NumberValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Cant use operator * with type "+type1.String(), node.Line, env)
		}
	} else if node.Operator == parser.OperatorDivide {

		if type1 == values.NumberType {
			if right.(values.NumberValue).Value == 0.0 {
				return e.Panic(values.ZeroDivisionError, "Division by zero", node.Line, env)
			}
			return values.NumberValue{Value: left.(values.NumberValue).Value / right.(values.NumberValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Cant use operator / with type "+type1.String(), node.Line, env)
		}
	}

	return e.Panic(values.RuntimeError, "Unknown operator", node.Line, env)

}

func (e Evaluator) EvaluateBinaryLogicExpression(node parser.BinaryLogicExpNode, env *environment.Environment) values.RuntimeValue {

	left := e.EvaluateExpression(node.Left, env)
	right := e.EvaluateExpression(node.Right, env)

	if left.GetType() == values.ErrorType {
		return left
	}
	if right.GetType() == values.ErrorType {
		return right
	}

	if node.Operator == parser.OperatorAnd {

		leftValue, err := e.EvaluateImplicitBoolConversion(left)

		if err != nil {
			return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
		}

		rightValue, err := e.EvaluateImplicitBoolConversion(right)

		if err != nil {
			return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
		}

		return values.BoolValue{Value: leftValue && rightValue}

	} else if node.Operator == parser.OperatorOr {

		leftValue, err := e.EvaluateImplicitBoolConversion(left)

		if err != nil {
			return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
		}

		rightValue, err := e.EvaluateImplicitBoolConversion(right)

		if err != nil {
			return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
		}

		return values.BoolValue{Value: leftValue || rightValue}

	}

	return values.ErrorValue{Value: "Unknown operator"}

}
func (e Evaluator) EvaluateBinaryComparisonExpression(node parser.BinaryComparisonExpNode, env *environment.Environment) values.RuntimeValue {

	left := e.EvaluateExpression(node.Left, env)
	right := e.EvaluateExpression(node.Right, env)

	type1 := left.GetType()
	type2 := right.GetType()

	if type1 == values.ErrorType {
		return left
	}
	if type2 == values.ErrorType {
		return right
	}

	equalTypes := type1 == type2

	if node.Operator == parser.OperatorEquals {

		if !equalTypes {
			return e.Panic(values.RuntimeError, "Type mismatch: "+type1.String()+" and "+type2.String(), node.Line, env)
		}

		if type1 == values.StringType {
			return values.BoolValue{Value: left.(values.StringValue).Value == right.(values.StringValue).Value}
		} else if type1 == values.NumberType {
			return values.BoolValue{Value: left.(values.NumberValue).Value == right.(values.NumberValue).Value}
		} else if type1 == values.BoolType {
			return values.BoolValue{Value: left.(values.BoolValue).Value == right.(values.BoolValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Cant use operator == with type "+type1.String(), node.Line, env)
		}
	} else if node.Operator == parser.OperatorGreaterThan {

		if !equalTypes {
			return e.Panic(values.RuntimeError, "Type mismatch: "+type1.String()+" and "+type2.String(), node.Line, env)
		}

		if type1 == values.NumberType {
			return values.BoolValue{Value: left.(values.NumberValue).Value > right.(values.NumberValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Operator > only can be used with numbers, not with type "+type1.String(), node.Line, env)
		}
	} else if node.Operator == parser.OperatorLessThan {

		if !equalTypes {
			return e.Panic(values.RuntimeError, "Type mismatch: "+type1.String()+" and "+type2.String(), node.Line, env)
		}

		if type1 == values.NumberType {
			return values.BoolValue{Value: left.(values.NumberValue).Value < right.(values.NumberValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Operator < only can be used with numbers, not with type "+type1.String(), node.Line, env)
		}
	} else if node.Operator == parser.OperatorLessOrEqThan {

		if !equalTypes {
			return e.Panic(values.RuntimeError, "Type mismatch: "+type1.String()+" and "+type2.String(), node.Line, env)
		}

		if type1 == values.NumberType {
			return values.BoolValue{Value: left.(values.NumberValue).Value <= right.(values.NumberValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Operator <= only can be used with numbers, not with type "+type1.String(), node.Line, env)
		}
	} else if node.Operator == parser.OperatorGreaterOrEqThan {

		if !equalTypes {
			return e.Panic(values.RuntimeError, "Type mismatch: "+type1.String()+" and "+type2.String(), node.Line, env)
		}

		if type1 == values.NumberType {
			return values.BoolValue{Value: left.(values.NumberValue).Value >= right.(values.NumberValue).Value}
		} else {
			return e.Panic(values.RuntimeError, "Operator >= only can be used with numbers, not with type "+type1.String(), node.Line, env)
		}
	}

	return values.ErrorValue{Value: "Unknown operator"}

}

func (e Evaluator) EvaluateUnaryExpression(node parser.UnaryExpNode, env *environment.Environment) values.RuntimeValue {

	exp := e.EvaluateExpression(node.Right, env)

	if node.Operator == "-" && exp.GetType() == values.NumberType {
		return values.NumberValue{Value: -exp.(values.NumberValue).Value}
	} else if node.Operator == "not" {
		res, err := e.EvaluateImplicitBoolConversion(exp)

		if err != nil {
			return e.Panic(values.RuntimeError, err.Error(), node.Line, env)
		}
		return values.BoolValue{Value: !res}
	} else {
		return values.ErrorValue{Value: "Unknown operator '" + node.Operator + "'"}
	}
}

func (e Evaluator) EvaluateImplicitBoolConversion(value values.RuntimeValue) (bool, error) {
	switch value.GetType() {
	case values.NumberType:
		return true, nil
	case values.StringType:
		return true, nil
	case values.BoolType:
		return value.(values.BoolValue).Value, nil
	case values.ArrayType:
		return true, nil
	default:
		return false, errors.New("Cannot convert " + value.GetType().String() + " to boolean")
	}
}

func (e Evaluator) ExecuteCallback(fn interface{}, args []interface{}) interface{} {

	fnValue := fn.(values.FunctionValue)
	fnEnv := fnValue.Environment.(*environment.Environment)

	for i, paramName := range fnValue.Parameters {
		if i+1 < len(args) {
			break
		}

		switch val := args[i].(type) {
		case *values.ObjectValue:
			fnEnv.Variables[paramName] = val
		case values.NumberValue:
			fnEnv.Variables[paramName] = val
		case values.StringValue:
			fnEnv.Variables[paramName] = val
		case values.BoolValue:
			fnEnv.Variables[paramName] = val
		case values.ArrayValue:
			fnEnv.Variables[paramName] = &val
		}
	}

	var result values.RuntimeValue

	for _, stmt := range fnValue.Body {
		result = e.EvaluateStmt(stmt, fnEnv)

		if result != nil && result.GetType() == values.ErrorType {
			return result
		}

		if result.GetType() == values.ReturnType {
			return result.(values.ReturnValue).Value
		}

	}
	return nil
}

func (e *Evaluator) AddToCallStack(line int, env *environment.Environment) {

	e.CallStack = append(e.CallStack, strconv.Itoa(line)+" "+env.ModuleName)
}

func (e *Evaluator) RemoveToCallStack() {
	e.CallStack = e.CallStack[:len(e.CallStack)-1]
}
