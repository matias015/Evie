package parser

import (
	"evie/lexer"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Parser context object
type ParserContext struct {
	AvoidStructInit bool
	Debug           bool
}

// Parser
type Parser struct {
	t       TokenIterator
	context ParserContext
}

func NewParser(tokens []lexer.Token) Parser {
	return Parser{t: TokenIterator{Items: tokens}, context: ParserContext{AvoidStructInit: false, Debug: false}}
}

func (p Parser) GetAST() []Stmt {
	// empty array
	ast := make([]Stmt, 0)

	for {

		if p.t.IsOutOfBounds() || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		parsed := p.ParseStmt()

		if parsed != nil {
			ast = append(ast, parsed)
		}
	}

	return ast
}

func (p *Parser) ParseStmt() Stmt {

	if p.t.Get().Kind == lexer.TOKEN_EOL {
		p.t.Eat()
	}

	if p.t.Get().Kind == lexer.TOKEN_EOF {
		return nil
	}

	token := p.t.Get()

	if token.Kind == lexer.TOKEN_VAR {
		return p.ParseVarDeclaration()
	} else if token.Kind == lexer.TOKEN_TRY {
		return p.ParseTryStmt()
	} else if token.Kind == lexer.TOKEN_IF {
		return p.ParseIfStmt()
	} else if token.Kind == lexer.TOKEN_IMPORT {
		return p.ParseImportStmt()
	} else if token.Kind == lexer.TOKEN_LOOP {
		return p.ParseLoopStmt()
	} else if token.Kind == lexer.TOKEN_FN && p.t.GetNext().Kind == lexer.TOKEN_IDENTIFIER {
		p.t.Eat()
		return p.ParseFunctionDeclaration()
	} else if token.Kind == lexer.TOKEN_FOR {
		return p.ParseForInStmt()
	} else if token.Kind == lexer.TOKEN_BREAK {
		return p.ParseBreakStmt()
	} else if token.Kind == lexer.TOKEN_CONTINUE {
		return p.ParseContinueStmt()
	} else if token.Kind == lexer.TOKEN_RETURN {
		return p.ParseReturnStmt()
	} else if token.Kind == lexer.TOKEN_IDENTIFIER && p.t.GetNext().Kind == lexer.TOKEN_LARROW {
		return p.ParseStructMethodDeclaration()
	} else if token.Kind == lexer.TOKEN_STRUCT {
		return p.ParseStructDeclaration()
	} else {

		return p.ParseExpressionStmt()
	}

}

func (p *Parser) ParseImportStmt() ImportNode {
	node := ImportNode{}
	node.Line = p.t.Eat().Line

	if p.t.Get().Kind != lexer.TOKEN_IDENTIFIER {
		Stop("Expected module name but found: " + lexer.GetTokenName(p.t.Get().Kind) + " in line " + fmt.Sprint(node.Line))
	}

	node.Path = p.t.Eat().Lexeme

	if p.t.Get().Kind == lexer.TOKEN_AS {
		p.t.Eat()
		if !p.t.HasNext() || p.t.GetNext().Kind != lexer.TOKEN_IDENTIFIER {
			Stop("Expected identifier after 'as' in line " + fmt.Sprint(node.Line))
		}
		node.Alias = p.t.Eat().Lexeme
	} else {
		node.Alias = node.Path

		if strings.ContainsRune(node.Path, '/') {
			splitted := strings.Split(node.Path, "/")
			node.Alias = splitted[len(splitted)-1]
		}
	}

	return node
}

func (p *Parser) ParseLoopStmt() LoopStmtNode {
	node := LoopStmtNode{}
	node.Line = p.t.Eat().Line

	node.Body = make([]Stmt, 0)

	if p.t.Get().Kind != lexer.TOKEN_LBRACE {
		Stop("Expected '{' in loop statement in line " + fmt.Sprint(p.t.Get().Line))
	}
	p.t.Eat() // {

	for {

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}
		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Get().Kind == lexer.TOKEN_EOL {
		p.t.Eat()
	}

	if p.t.Get().Kind == lexer.TOKEN_RBRACE {
		p.t.Eat()
	} else {
		Stop("Expected '}' in loop statement in line " + fmt.Sprint(p.t.Get().Line))
	}

	if len(node.Body) == 0 {
		Stop("Empty loop statement in line " + fmt.Sprint(p.t.Get().Line))
	}
	return node
}

func (p *Parser) ParseTryStmt() TryCatchNode {
	node := TryCatchNode{}
	node.Line = p.t.Eat().Line

	node.Body = make([]Stmt, 0)

	if p.t.Get().Kind != lexer.TOKEN_LBRACE {
		Stop("Expected '{' in try statement in line " + fmt.Sprint(p.t.Get().Line))
	}

	p.t.Eat() // {

	for {

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}
		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Get().Kind != lexer.TOKEN_RBRACE {
		Stop("Expected '}' in try statement in line " + fmt.Sprint(p.t.Get().Line))
	}
	p.t.Eat()

	if p.t.Eat().Kind != lexer.TOKEN_CATCH {
		Stop("expected catch")
	} else {
		p.t.Eat()
	}

	node.Catch = make([]Stmt, 0)

	for {

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}
		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		node.Catch = append(node.Catch, p.ParseStmt())
	}
	if p.t.Get().Kind != lexer.TOKEN_RBRACE {
		Stop("Expected '}' in catch statement in line " + fmt.Sprint(p.t.Get().Line))
	}
	p.t.Eat()

	if p.t.Get().Kind == lexer.TOKEN_FINALLY {
		p.t.Eat()

		node.Finally = make([]Stmt, 0)

		if p.t.Get().Kind == lexer.TOKEN_LBRACE {
			p.t.Eat()
		} else {
			Stop("Expected '{' in finally statement in line " + fmt.Sprint(p.t.Get().Line))
		}

		for {
			if p.t.Get().Kind == lexer.TOKEN_EOL {
				p.t.Eat()
				continue
			}

			if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
				break
			}

			node.Finally = append(node.Finally, p.ParseStmt())
		}

		if p.t.Get().Kind == lexer.TOKEN_RBRACE {
			p.t.Eat()
		} else {
			Stop("Expected '}' in finally statement in line " + fmt.Sprint(p.t.Get().Line))
		}

	}

	return node
}

func (p *Parser) ParseReturnStmt() ReturnNode {
	node := ReturnNode{}
	node.Line = p.t.Eat().Line
	node.Right = p.ParseExp()

	return node
}
func (p *Parser) ParseBreakStmt() BreakNode {
	return BreakNode{Line: p.t.Eat().Line}
}

func (p *Parser) ParseContinueStmt() ContinueNode {
	return ContinueNode{Line: p.t.Eat().Line}
}

func (p *Parser) ParseForInStmt() ForInSatementNode {
	node := ForInSatementNode{}

	node.Body = make([]Stmt, 0)

	node.Line = p.t.Eat().Line

	if p.t.Get().Kind != lexer.TOKEN_IDENTIFIER {
		Stop("Expected at least one identifier after 'for' keyword in for-in statement in line " + fmt.Sprint(p.t.Get().Line))
	}

	firstVar := p.t.Eat().Lexeme

	var secondVar string = ""

	if p.t.Get().Kind == lexer.TOKEN_COMMA {
		p.t.Eat()
		secondVar = p.t.Eat().Lexeme
	}

	if secondVar != "" {
		node.IndexVarName = firstVar
		node.LocalVarName = secondVar
	} else {
		node.LocalVarName = firstVar
	}

	if p.t.Get().Kind == lexer.TOKEN_IN {
		p.t.Eat()
	} else {
		Stop("Expected 'in' keyword in for in statement in line " + fmt.Sprint(p.t.Get().Line))
	}

	p.context.AvoidStructInit = true
	iterator := p.ParseExp()
	p.context.AvoidStructInit = false

	node.Iterator = iterator

	if p.t.Eat().Kind != lexer.TOKEN_LBRACE {
		Stop("Expected '{' in for in statement in line " + fmt.Sprint(p.t.Get().Line))
	}

	for {

		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Get().Kind != lexer.TOKEN_RBRACE {
		Stop("Expected '}' in for in statement in line " + fmt.Sprint(p.t.Get().Line))
	} else {
		p.t.Eat()
	}

	if len(node.Body) == 0 {
		Stop("Empty loop statement in line " + fmt.Sprint(p.t.Get().Line))
	}
	return node
}

func (p *Parser) ParseStructMethodDeclaration() StructMethodDeclarationNode {
	node := StructMethodDeclarationNode{}

	structNameToken := p.t.Eat()
	node.Line = structNameToken.Line

	node.Struct = structNameToken.Lexeme

	p.t.Eat() // ->

	node.Function = p.ParseFunctionDeclaration()

	return node
}

func (p *Parser) ParseStructDeclaration() StructDeclarationNode {
	line := p.t.Eat().Line // struct keyword

	var node StructDeclarationNode = StructDeclarationNode{}

	node.Line = line

	// struct name

	if p.t.Get().Kind != lexer.TOKEN_IDENTIFIER {
		Stop("Expected identifier after 'struct' keyword in struct declaration in line " + fmt.Sprint(p.t.Get().Line))
	}

	node.Name = p.t.Eat().Lexeme
	node.Properties = make([]string, 0)

	if p.t.Get().Kind != lexer.TOKEN_LBRACE {
		Stop("Expected '{' in struct declaration in line " + fmt.Sprint(node.Line))
	}

	p.t.Eat()

	for {
		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}

		if p.t.Get().Kind != lexer.TOKEN_IDENTIFIER {
			Stop("Expected identifier in struct declaration but found: " + lexer.GetTokenName(p.t.Get().Kind) + " in line " + fmt.Sprint(node.Line))
		}

		node.Properties = append(node.Properties, p.t.Eat().Lexeme)

		if p.t.Get().Kind == lexer.TOKEN_COMMA {
			p.t.Eat()
		}
	}

	if p.t.Get().Kind != lexer.TOKEN_RBRACE {
		Stop("Expected '}' in struct declaration in line " + fmt.Sprint(node.Line))
	}
	p.t.Eat()

	return node
}

func (p *Parser) ParseFunctionDeclaration() FunctionDeclarationNode {

	var node FunctionDeclarationNode = FunctionDeclarationNode{}
	node.Body = make([]Stmt, 0)

	if p.t.Get().Kind == lexer.TOKEN_IDENTIFIER {
		node.Name = p.t.Get().Lexeme
		node.Line = p.t.Eat().Line
	} else if p.t.Get().Kind == lexer.TOKEN_LPAR {
		node.Name = ""
		node.Line = p.t.Get().Line
	} else {
		Stop("Expected identifier or '(' in function declaration in line " + fmt.Sprint(p.t.Get().Line))
	}

	args := p.ParseArgs()

	for _, arg := range args {
		node.Parameters = append(node.Parameters, arg.(IdentifierNode).Value)
	}

	if p.t.Get().Kind != lexer.TOKEN_LBRACE {
		Stop("Expected '{' in function declaration in line " + fmt.Sprint(p.t.Get().Line))
	}

	p.t.Eat() // open brace

	for {

		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}
		node.Body = append(node.Body, p.ParseStmt())

	}

	if p.t.Get().Kind == lexer.TOKEN_EOL {
		p.t.Eat()
	}

	if p.t.Get().Kind != lexer.TOKEN_RBRACE {
		Stop("Expected '}' in function declaration in line " + fmt.Sprint(p.t.Get().Line))
	} else {
		p.t.Eat()
	}

	return node
}

func (p *Parser) ParseIfStmt() IfStatementNode {
	// If token
	line := p.t.Eat().Line

	// Node
	var node IfStatementNode = IfStatementNode{}

	node.Line = line

	// Parse condition
	p.context.AvoidStructInit = true
	node.Condition = p.ParseExp()
	p.context.AvoidStructInit = false

	p.t.Eat() // open brace

	for {
		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Get().Kind == lexer.TOKEN_EOL {
		p.t.Eat()
	}
	p.t.Eat() // close brace

	// ELSE IF

	node.ElseIf = make([]IfStatementNode, 0)
	for {
		if !(p.t.Get().Kind == lexer.TOKEN_ELSE && p.t.GetNext().Kind == lexer.TOKEN_IF) {
			break
		}

		p.t.Eat() // else
		p.t.Eat() // if

		var elseifnode IfStatementNode = IfStatementNode{}
		elseifnode.Body = make([]Stmt, 0)

		elseifnode.Line = line

		p.context.AvoidStructInit = true
		elseifnode.Condition = p.ParseExp()
		p.context.AvoidStructInit = false

		p.t.Eat() // open brace
		for {
			if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
				break
			}

			if p.t.Get().Kind == lexer.TOKEN_EOL {
				p.t.Eat()
				continue
			}

			elseifnode.Body = append(elseifnode.Body, p.ParseStmt())
		}

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
		}
		if p.t.Get().Kind == lexer.TOKEN_RBRACE {
			p.t.Eat()
		}

		node.ElseIf = append(node.ElseIf, elseifnode)
	}

	if p.t.Get().Kind == lexer.TOKEN_ELSE {
		node.ElseBody = make([]Stmt, 0)
		p.t.Eat()
		p.t.Eat()
		for {
			if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
				break
			}

			if p.t.Get().Kind == lexer.TOKEN_EOL {
				p.t.Eat()
				continue
			}

			node.ElseBody = append(node.ElseBody, p.ParseStmt())
		}

		p.t.Eat()
		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
		}
	}

	return node
}

func (p *Parser) ParseExpressionStmt() ExpressionStmtNode {
	token := p.t.Get()
	val := p.ParseExp()
	if val == nil {
		fmt.Println("Bad expression after token " + token.Lexeme + " line " + strconv.Itoa(token.Line))
		os.Exit(1)
	}
	return ExpressionStmtNode{Expression: val}
}

func (p *Parser) ParseVarDeclaration() Stmt {
	line := p.t.Eat().Line

	node := VarDeclarationNode{}

	node.Line = line

	identifier := p.t.Eat()

	node.Left = IdentifierNode{Value: identifier.Lexeme}

	if p.t.Get().Kind == lexer.TOKEN_EOL {
		node.Right = NothingNode{Line: line}
	} else {
		operator := p.t.Eat()
		node.Operator = operator.Lexeme

		node.Right = p.ParseExp()
	}

	return node
}

func (p *Parser) ParseExp() Exp {

	return p.ParseAnonFnExp()

}

func (p *Parser) ParseAnonFnExp() Exp {

	if p.t.Get().Lexeme != "fn" {
		return p.ParseAssignmentExp()
	}

	p.t.Eat()

	var node AnonFunctionDeclarationNode = AnonFunctionDeclarationNode{}
	node.Body = make([]Stmt, 0)

	args := p.ParseArgs()

	for _, arg := range args {
		node.Parameters = append(node.Parameters, arg.(IdentifierNode).Value)
	}

	if p.t.Get().Kind != lexer.TOKEN_LBRACE {
		Stop("Expected '{' in anon function declaration in line " + fmt.Sprint(p.t.Get().Line))
	}

	p.t.Eat() // open brace

	for {

		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}
		node.Body = append(node.Body, p.ParseStmt())

	}

	if p.t.Get().Kind == lexer.TOKEN_EOL {
		p.t.Eat()
	}

	if p.t.Get().Kind != lexer.TOKEN_RBRACE {
		Stop("Expected '}' in function declaration in line " + fmt.Sprint(p.t.Get().Line))
	} else {
		p.t.Eat()
	}

	return node

}

func (p *Parser) ParseAssignmentExp() Exp {
	left := p.ParseTernaryExp()

	if p.t.Get().Lexeme == "=" {

		operator := p.t.Eat()
		right := p.ParseTernaryExp()
		n := AssignmentNode{}
		n.Line = operator.Line
		n.Left = left
		n.Operator = operator.Lexeme
		n.Right = right
		left = n
	}
	return left
}

func (p *Parser) ParseTernaryExp() Exp {

	left := p.ParseDictionaryInitialization()

	if p.t.Get().Lexeme == "?" {

		n := TernaryExpNode{}
		n.Line = p.t.Eat().Line

		n.Condition = left

		n.Left = p.ParseDictionaryInitialization()

		if p.t.Get().Lexeme != ":" {
			Stop("Expected ':'")
		}
		p.t.Eat()

		n.Right = p.ParseDictionaryInitialization()

		left = n
	}

	return left
}

func (p *Parser) ParseDictionaryInitialization() Exp {

	if p.t.Get().Lexeme != "{" || p.context.AvoidStructInit == true {
		return p.ParseBinaryExp()
	}

	line := p.t.Eat().Line

	node := DictionaryExpNode{}

	node.Line = line

	node.Value = make(map[string]Exp, 0)

	for {

		if p.t.Get().Kind == lexer.TOKEN_RBRACE || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		if p.t.Get().Lexeme == "," {
			p.t.Eat()
			continue
		}

		// ignoramos fin de linea
		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}

		// Agregamos la llave
		key := p.t.Eat()

		// Luego deberia haber un :
		// ex -> {a: 2, b: 3}

		if p.t.Eat().Lexeme != ":" {
			Stop("Expected ':' in line " + fmt.Sprint(p.t.Get().Line) + " after key " + key.Lexeme)
		}

		value := p.ParseExp()

		node.Value[key.Lexeme] = value

	}

	if p.t.Get().Kind == lexer.TOKEN_EOL {
		p.t.Eat()
	}
	p.t.Eat()

	return node
}

func (p *Parser) ParseBinaryExp() Exp {
	return p.parseLogicOrExpression()
}

func (p *Parser) parseLogicOrExpression() Exp {
	left := p.parseLogicAndExpression()

	for p.t.Get().Lexeme == "or" {
		op := p.t.Eat()
		n := BinaryExpNode{}
		n.Line = op.Line
		n.Left = left
		n.Operator = op.Lexeme
		n.Right = p.parseLogicAndExpression()
		left = n
	}

	return left
}

func (p *Parser) parseLogicAndExpression() Exp {
	left := p.parseComparisonExp()

	for p.t.Get().Lexeme == "and" {
		op := p.t.Eat()
		n := BinaryExpNode{}
		n.Line = op.Line
		n.Left = left
		n.Operator = op.Lexeme
		n.Right = p.parseComparisonExp()
		left = n
	}

	return left

}

func (p *Parser) parseComparisonExp() Exp {
	left := p.parseAdditiveExp()

	for p.t.Get().Lexeme == "==" || p.t.Get().Lexeme == "!=" || p.t.Get().Lexeme == ">" || p.t.Get().Lexeme == "<" || p.t.Get().Lexeme == ">=" || p.t.Get().Lexeme == "<=" {
		op := p.t.Eat()
		n := BinaryExpNode{}
		n.Line = op.Line
		n.Left = left
		n.Operator = op.Lexeme
		n.Right = p.parseAdditiveExp()
		left = n
	}

	return left
}

func (p *Parser) parseAdditiveExp() Exp {
	left := p.parseObjectInitExp()

	for p.t.Get().Lexeme == "+" || p.t.Get().Lexeme == "-" {
		op := p.t.Eat()
		n := BinaryExpNode{}
		n.Left = left
		n.Line = op.Line
		n.Operator = op.Lexeme
		n.Right = p.parseObjectInitExp()
		left = n
	}

	return left

}

func (p *Parser) parseObjectInitExp() Exp {
	left := p.parseMultiplicativeExp()

	for p.t.Get().Lexeme == "{" && p.context.AvoidStructInit == false {

		node := ObjectInitExpNode{}
		node.Line = p.t.Get().Line

		p.context.AvoidStructInit = true
		p.context.AvoidStructInit = false

		node.Struct = left
		node.Value = p.ParseDictionaryInitialization().(DictionaryExpNode)

		left = node
	}

	return left

}

func (p *Parser) parseMultiplicativeExp() Exp {
	left := p.ParseMemberExp()

	for p.t.Get().Lexeme == "*" || p.t.Get().Lexeme == "/" {
		op := p.t.Eat()
		n := BinaryExpNode{}
		n.Left = left
		n.Line = op.Line
		n.Operator = op.Lexeme
		n.Right = p.ParseMemberExp()
		left = n
	}

	return left
}
func (p *Parser) ParseMemberExp() Exp {

	left := p.ParseCallMemberExp()

	for p.t.Get().Lexeme == "." {

		line := p.t.Eat().Line

		n := MemberExpNode{}
		n.Line = line
		n.Left = left
		n.Member = IdentifierNode{Value: p.t.Eat().Lexeme}
		left = n

		if p.t.Get().Lexeme == "(" {
			left = p.ParseCallExpr(left)
		}
	}

	return left

}

func (p *Parser) ParseCallMemberExp() Exp {
	member := p.ParseIndexAccessExp()

	if p.t.Get().Lexeme == "(" {
		return p.ParseCallExpr(member)
	}

	return member
}

func (p *Parser) ParseIndexAccessExp() Exp {

	left := p.parseUnaryExp()

	for p.t.Get().Lexeme == "[" {

		line := p.t.Eat().Line

		if p.t.Get().Lexeme == ":" {
			p.t.Eat()

			if p.t.Get().Lexeme == "]" {
				Stop("Empty slice expression at line " + fmt.Sprint(line))
			}

			sliceNode := SliceExpNode{}
			sliceNode.Line = line
			sliceNode.Left = left
			sliceNode.From = NumberNode{Value: "0"}
			sliceNode.To = p.ParseExp()
			left = sliceNode

			if p.t.Get().Lexeme != "]" {
				Stop("Expected ']'")
			}

			p.t.Eat()

			return left
		}

		index := p.ParseExp()

		n := IndexAccessExpNode{}
		n.Line = line
		n.Left = left
		n.Index = index

		if p.t.Get().Lexeme == ":" {
			p.t.Eat()
			sliceNode := SliceExpNode{}
			sliceNode.Line = line
			sliceNode.Left = left
			sliceNode.From = index
			if p.t.Get().Lexeme == "]" {
				sliceNode.To = nil
			} else {
				sliceNode.To = p.ParseExp()
			}
			left = sliceNode
		} else {
			left = n
		}

		if p.t.Get().Lexeme == "]" {
			p.t.Eat()
		} else {
			Stop("Expected ']' at line " + fmt.Sprint(line))
		}
	}

	return left

}

func (p *Parser) parseUnaryExp() Exp {

	if p.t.Get().Lexeme == "-" && p.t.Get().Kind == lexer.TOKEN_OPERATOR {
		op := p.t.Eat()
		n := UnaryExpNode{}
		n.Line = op.Line
		n.Operator = op.Lexeme
		n.Right = p.parseUnaryExp()
		return n
	} else if p.t.Get().Lexeme == "not" {
		op := p.t.Eat()
		n := UnaryExpNode{}
		n.Line = op.Line
		n.Operator = op.Lexeme
		n.Right = p.parseUnaryExp()
		return n
	}
	return p.parsePrimaryExp()
}

func (p *Parser) parsePrimaryExp() Exp {

	token := p.t.Eat()

	if token.Kind == lexer.TOKEN_IDENTIFIER {
		return IdentifierNode{Value: token.Lexeme, Line: token.Line}
	} else if token.Kind == lexer.TOKEN_NUMBER {
		return NumberNode{Value: token.Lexeme, Line: token.Line}
	} else if token.Kind == lexer.TOKEN_NOTHING {
		return NothingNode{Line: token.Line}
	} else if token.Kind == lexer.TOKEN_STRING {
		return StringNode{Value: token.Lexeme, Line: token.Line}
	} else if token.Kind == lexer.TOKEN_LBRACKET {
		return p.ParseArrayInitializationExp()
	} else if token.Kind == lexer.TOKEN_BOOLEAN {
		var value bool
		if token.Lexeme == "true" {
			value = true
		} else if token.Lexeme == "false" {
			value = false
		}
		return BooleanNode{Value: value, Line: token.Line}
	} else if token.Kind == lexer.TOKEN_LPAR {
		v := p.ParseExp()
		p.t.Eat()
		return v
	} else if token.Kind == lexer.TOKEN_EOF {
		Stop("Unexpected end of file in line " + fmt.Sprint(token.Line))
		return nil
	} else if token.Kind == lexer.TOKEN_EOL {
		Stop("Unexpected line break in line " + fmt.Sprint(token.Line))
		return nil
	} else {
		Stop("Unexpected token: " + token.Lexeme + " in line " + fmt.Sprint(token.Line))
		return nil
	}
}

func (p *Parser) ParseCallExpr(member Exp) Exp {

	node := CallExpNode{}
	node.Line = p.t.Get().Line

	if member.ExpType() != "IdentifierNode" && member.ExpType() != "MemberExpNode" {
		Stop("Expected identifier in call expression in line " + fmt.Sprint(p.t.Get().Line))
	}

	node.Name = member

	args := p.ParseArgs()

	node.Args = args

	return node
}

func (p *Parser) ParseArgs() []Exp {

	if p.t.Get().Kind != lexer.TOKEN_LPAR {
		Stop("Expected '(' after function arguments in line " + fmt.Sprint(p.t.Get().Line))
	}

	p.t.Eat()

	var args []Exp

	if p.t.Get().Lexeme == ")" {
		p.t.Eat()
		return args
	}

	args = p.ParseArgumentsList()

	return args

}

// ParseArgumentsList parses a list of expressions separated by commas.
// It returns a slice of expressions representing the parsed arguments.
func (p *Parser) ParseArgumentsList() []Exp {
	var args []Exp

	p.context.AvoidStructInit = true

	args = append(args, p.ParseExp())

	for p.t.Get().Lexeme == "," {
		p.t.Eat()
		args = append(args, p.ParseExp())
	}

	if p.t.Get().Lexeme != ")" {
		Stop("Expected ')' after function arguments in line " + fmt.Sprint(p.t.Get().Line))
	}
	p.context.AvoidStructInit = false
	p.t.Eat()

	return args
}

func (p *Parser) ParseArrayInitializationExp() Exp {
	node := ArrayExpNode{}
	node.Line = p.t.Get().Line
	node.Value = make([]Exp, 0)

	for {

		if p.t.Get().Kind == lexer.TOKEN_RBRACKET || p.t.Get().Kind == lexer.TOKEN_EOF {
			break
		}

		if p.t.Get().Kind == lexer.TOKEN_COMMA {
			p.t.Eat()
			continue
		}

		if p.t.Get().Kind == lexer.TOKEN_EOL {
			p.t.Eat()
			continue
		}
		// "a"
		node.Value = append(node.Value, p.ParseExp())

	}

	if p.t.Get().Kind != lexer.TOKEN_RBRACKET {
		Stop("Expected ']' in line " + fmt.Sprint(p.t.Get().Line))
	} else {
		p.t.Eat()
	}

	return node
}

func Stop(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (p *Parser) Debug(text string) {
	if p.context.Debug {
		fmt.Println(text)
	}
}
