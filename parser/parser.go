package parser

import (
	"evie/lexer"
	"fmt"
	"os"
	"strings"
)

// Parser context object
type ParserContext struct {
	AvoidStructInit bool
}

// Parser
type Parser struct {
	t       TokenIterator
	context ParserContext
}

func NewParser(tokens []lexer.Token) Parser {
	return Parser{t: TokenIterator{Items: tokens}, context: ParserContext{AvoidStructInit: false}}
}

func (p Parser) GetAST() []Stmt {
	// empty array
	ast := make([]Stmt, 0)

	for {

		if p.t.IsOutOfBounds() || p.t.Get().Kind == "eof" {
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

	if p.t.Get().Kind == "eol" {
		p.t.Eat()
	}

	if p.t.Get().Kind == "eof" {
		return nil
	}

	token := p.t.Get()

	if token.Kind == "var" {
		return p.ParseVarDeclaration()
	} else if token.Kind == "try" {
		return p.ParseTryStmt()
	} else if token.Kind == "if" {
		return p.ParseIfStmt()
	} else if token.Kind == "import" {
		return p.ParseImportStmt()
	} else if token.Kind == "loop" {
		return p.ParseLoopStmt()
	} else if token.Kind == "fn" {
		p.t.Eat()
		return p.ParseFunctionDeclaration()
	} else if token.Kind == "for" {
		return p.ParseForInStmt()
	} else if token.Kind == "break" {
		return p.ParseBreakStmt()
	} else if token.Kind == "continue" {
		return p.ParseContinueStmt()
	} else if token.Kind == "return" {
		return p.ParseReturnStmt()
	} else if token.Kind == "identifier" && p.t.GetNext().Kind == "arrowleft" {
		return p.ParseStructMethodDeclaration()
	} else if token.Kind == "struct" {
		return p.ParseStructDeclaration()
	} else {
		return p.ParseExpressionStmt()
	}

}

func (p *Parser) ParseImportStmt() ImportNode {
	node := ImportNode{}
	node.Line = p.t.Eat().Line
	node.Path = p.t.Eat().Lexeme

	if p.t.Get().Kind == "as" {
		p.t.Eat()
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

	p.t.Eat() // {

	for {

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}
		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Get().Kind == "rbrace" {
		p.t.Eat()
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

	p.t.Eat() // {

	for {

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}
		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Get().Kind == "rbrace" {
		p.t.Eat()
	}

	if p.t.Eat().Kind != "catch" {
		Stop("expected catch")
	} else {
		p.t.Eat()
	}

	node.Catch = make([]Stmt, 0)

	for {

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}
		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		node.Catch = append(node.Catch, p.ParseStmt())
	}
	if p.t.Get().Kind == "rbrace" {
		p.t.Eat()
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

	firstVar := p.t.Eat().Lexeme

	var secondVar string = ""

	if p.t.Get().Kind == "comma" {
		p.t.Eat()
		secondVar = p.t.Eat().Lexeme
	}

	if secondVar != "" {
		node.IndexVarName = firstVar
		node.LocalVarName = secondVar
	} else {
		node.LocalVarName = firstVar
	}

	if p.t.Get().Kind == "in" {
		p.t.Eat()
	} else {
		Stop("Expected 'in' keyword in for in statement in line " + fmt.Sprint(p.t.Get().Line))
	}

	p.context.AvoidStructInit = true
	iterator := p.ParseExp()
	p.context.AvoidStructInit = false

	node.Iterator = iterator

	if p.t.Eat().Kind != "lbrace" {
		Stop("Expected '{' in for in statement in line " + fmt.Sprint(p.t.Get().Line))
	}

	for {

		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Eat().Kind != "rbrace" {
		Stop("Expected '}' in for in statement in line " + fmt.Sprint(p.t.Get().Line))
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
	node.Name = p.t.Eat().Lexeme
	node.Properties = make([]string, 0)

	p.t.Eat()

	for {
		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}

		node.Properties = append(node.Properties, p.t.Eat().Lexeme)

		if p.t.Get().Kind == "comma" {
			p.t.Eat()
		}
	}

	p.t.Eat()

	if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" || p.t.Get().Kind == "eol" {
		p.t.Eat()
	}

	return node
}

func (p *Parser) ParseFunctionDeclaration() FunctionDeclarationNode {

	var node FunctionDeclarationNode = FunctionDeclarationNode{}
	node.Body = make([]Stmt, 0)
	node.Name = p.t.Get().Lexeme
	node.Line = p.t.Eat().Line

	args := p.ParseArgs()

	for _, arg := range args {
		node.Parameters = append(node.Parameters, arg.(IdentifierNode).Value)
	}

	p.t.Eat() // open brace

	for {
		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	p.t.Eat()

	if p.t.Get().Kind == "eol" {
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
		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}

		node.Body = append(node.Body, p.ParseStmt())
	}

	if p.t.Get().Kind == "eol" {
		p.t.Eat()
	}
	p.t.Eat() // close brace

	// ELSE IF

	node.ElseIf = make([]IfStatementNode, 0)

	for {
		if !(p.t.Get().Kind == "else" && p.t.GetNext().Lexeme == "if") {
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
			if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
				break
			}

			if p.t.Get().Kind == "eol" {
				p.t.Eat()
				continue
			}

			elseifnode.Body = append(elseifnode.Body, p.ParseStmt())
		}

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
		}
		if p.t.Get().Kind == "rbrace" {
			p.t.Eat()
		}

		node.ElseIf = append(node.ElseIf, elseifnode)
	}

	if p.t.Get().Kind == "else" {
		node.ElseBody = make([]Stmt, 0)
		p.t.Eat()
		p.t.Eat()
		for {
			if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
				break
			}

			if p.t.Get().Kind == "eol" {
				p.t.Eat()
				continue
			}

			node.ElseBody = append(node.ElseBody, p.ParseStmt())
		}

		p.t.Eat()
		if p.t.Get().Kind == "eol" {
			p.t.Eat()
		}
	}

	return node
}

func (p *Parser) ParseExpressionStmt() ExpressionStmtNode {
	return ExpressionStmtNode{Expression: p.ParseExp()}
}

func (p *Parser) ParseVarDeclaration() Stmt {
	line := p.t.Eat().Line

	node := VarDeclarationNode{}

	node.Line = line

	identifier := p.t.Eat()
	operator := p.t.Eat()
	expression := p.ParseExp()

	node.Left = IdentifierNode{Value: identifier.Lexeme}
	node.Operator = operator.Lexeme
	node.Right = expression

	return node
}

func (p *Parser) ParseExp() Exp {

	return p.ParseAssignmentExp()

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

	if p.t.Get().Lexeme != "{" {
		return p.ParseBinaryExp()
	}

	line := p.t.Eat().Line

	node := DictionaryExpNode{}

	node.Line = line

	node.Value = make(map[string]Exp, 0)

	for {

		if p.t.Get().Kind == "rbrace" || p.t.Get().Kind == "eof" {
			break
		}

		if p.t.Get().Lexeme == "," {
			p.t.Eat()
			continue
		}

		// ignoramos fin de linea
		if p.t.Get().Kind == "eol" {
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
	p.t.Eat()

	if p.t.Get().Kind == "eol" || p.t.Get().Kind == "rbrace" {
		p.t.Eat()
	}

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
		node.Struct = left.(IdentifierNode).Value
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
			sliceNode := SliceExpNode{}
			sliceNode.Line = line
			sliceNode.Left = left
			sliceNode.From = nil
			sliceNode.To = p.ParseExp()
			left = sliceNode
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

		p.t.Eat()

	}

	return left

}

func (p *Parser) parseUnaryExp() Exp {
	if p.t.Get().Lexeme == "-" && p.t.Get().Kind == "operator" {
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

	if token.Kind == "identifier" {
		return IdentifierNode{Value: token.Lexeme, Line: token.Line}
	} else if token.Kind == "number" {
		return NumberNode{Value: token.Lexeme, Line: token.Line}
	} else if token.Kind == "string" {
		return StringNode{Value: token.Lexeme, Line: token.Line}
	} else if token.Kind == "lbracket" {
		return p.ParseArrayInitializationExp()
	} else if token.Kind == "boolean" {
		var value bool
		if token.Lexeme == "true" {
			value = true
		} else if token.Lexeme == "false" {
			value = false
		}
		return BooleanNode{Value: value, Line: token.Line}
	} else if token.Kind == "lpar" {
		v := p.ParseExp()
		p.t.Eat()
		return v
	} else {
		return nil
	}
}

func (p *Parser) ParseCallExpr(member Exp) Exp {

	node := CallExpNode{}
	node.Line = p.t.Get().Line

	node.Name = member

	args := p.ParseArgs()

	node.Args = args

	return node
}

func (p *Parser) ParseArgs() []Exp {
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

	args = append(args, p.ParseExp())

	for p.t.Get().Lexeme == "," {
		p.t.Eat()
		args = append(args, p.ParseExp())
	}

	p.t.Eat()

	return args
}

func (p *Parser) ParseArrayInitializationExp() Exp {
	node := ArrayExpNode{}
	node.Line = p.t.Get().Line
	node.Value = make([]Exp, 0)

	for {
		if p.t.Get().Kind == "rbracket" || p.t.Get().Kind == "eof" {
			p.t.Eat()
			break
		}

		if p.t.Get().Kind == "comma" {
			p.t.Eat()
			continue
		}

		if p.t.Get().Kind == "eol" {
			p.t.Eat()
			continue
		}

		node.Value = append(node.Value, p.ParseExp())

	}

	return node
}

func Stop(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
