package parser

// INTERFACES
type Stmt interface {
	StmtType() string
}

type Exp interface {
	ExpType() string
}

// EXPRESIONES
type ExpressionStmtNode struct {
	Expression Exp
}

func (e ExpressionStmtNode) StmtType() string { return "ExpressionStmtNode" }

type NumberNode struct {
	Value string
	Line  int
}

func (n NumberNode) ExpType() string { return "NumberNode" }

type StringNode struct {
	Value string
	Line  int
}

func (n StringNode) ExpType() string { return "StringNode" }

type BooleanNode struct {
	Value bool
	Line  int
}

func (n BooleanNode) ExpType() string { return "BooleanNode" }

type IdentifierNode struct {
	Value string
	Line  int
}

func (n IdentifierNode) ExpType() string { return "IdentifierNode" }

type NothingNode struct {
	Line int
}

func (n NothingNode) ExpType() string { return "NothingNode" }

type AssignmentNode struct {
	Left     Exp
	Operator string
	Right    Exp
	Line     int
}

func (n AssignmentNode) ExpType() string { return "AssignmentNode" }

type BinaryExpNode struct {
	Left     Exp
	Operator string
	Right    Exp
	Line     int
}

func (n BinaryExpNode) ExpType() string { return "BinaryExpNode" }

type UnaryExpNode struct {
	Operator string
	Right    Exp
	Line     int
}

func (n UnaryExpNode) ExpType() string { return "UnaryExpNode" }

type CallExpNode struct {
	Name Exp
	Args []Exp
	Line int
}

func (n CallExpNode) ExpType() string { return "CallExpNode" }

type ArrayExpNode struct {
	Value []Exp
	Line  int
}

func (n ArrayExpNode) ExpType() string { return "ArrayExpNode" }

type IndexAccessExpNode struct {
	Left  Exp
	Index Exp
	Line  int
}

func (n IndexAccessExpNode) ExpType() string { return "IndexAccessExpNode" }

type DictionaryExpNode struct {
	Value map[string]Exp
	Line  int
}

func (n DictionaryExpNode) ExpType() string { return "DictionaryExpNode" }

// Struct Initialization, NOT ANY OBJECT LIKE JS
type ObjectInitExpNode struct {
	Struct Exp
	Value  DictionaryExpNode
	Line   int
}

func (n ObjectInitExpNode) ExpType() string { return "ObjectInitExpNode" }

// Struct Initialization, NOT ANY OBJECT LIKE JS
type MemberExpNode struct {
	Left   Exp
	Member IdentifierNode
	Line   int
}

func (n MemberExpNode) ExpType() string { return "MemberExpNode" }

// slice

type SliceExpNode struct {
	Left Exp
	From Exp
	To   Exp
	Line int
}

func (n SliceExpNode) ExpType() string { return "SliceExpNode" }

type TernaryExpNode struct {
	Condition Exp
	Left      Exp
	Right     Exp
	Line      int
}

func (n TernaryExpNode) ExpType() string { return "TernaryExpNode" }

// STATEMENTS

type VarDeclarationNode struct {
	Left     IdentifierNode
	Operator string
	Right    Exp
	Line     int
}

func (n VarDeclarationNode) StmtType() string { return "VarDeclarationNode" }

type IfStatementNode struct {
	Condition Exp
	Body      []Stmt
	ElseIf    []IfStatementNode
	ElseBody  []Stmt
	Line      int
}

func (n IfStatementNode) StmtType() string { return "IfStatementNode" }

type FunctionDeclarationNode struct {
	Name       string
	Body       []Stmt
	Parameters []string
	Line       int
}

func (n FunctionDeclarationNode) StmtType() string { return "FunctionDeclarationNode" }

type StructDeclarationNode struct {
	Name       string
	Properties []string
	Line       int
}

func (n StructDeclarationNode) StmtType() string { return "StructDeclarationNode" }

type StructMethodDeclarationNode struct {
	Struct   string
	Function FunctionDeclarationNode
	Line     int
}

func (n StructMethodDeclarationNode) StmtType() string { return "StructMethodDeclarationNode" }

type ForInSatementNode struct {
	Iterator     Exp
	Body         []Stmt
	IndexVarName string
	LocalVarName string
	Line         int
}

func (n ForInSatementNode) StmtType() string { return "ForInSatementNode" }

type BreakNode struct {
	Line int
}

func (n BreakNode) StmtType() string { return "BreakNode" }

type ContinueNode struct {
	Line int
}

func (n ContinueNode) StmtType() string { return "ContinueNode" }

type ReturnNode struct {
	Right Exp
	Line  int
}

func (n ReturnNode) StmtType() string { return "ReturnNode" }

type TryCatchNode struct {
	Body    []Stmt
	Catch   []Stmt
	Finally []Stmt
	Line    int
}

func (n TryCatchNode) StmtType() string { return "TryCatchNode" }

type LoopStmtNode struct {
	Body []Stmt
	Line int
}

func (n LoopStmtNode) StmtType() string { return "LoopStmtNode" }

type ImportNode struct {
	Line  int
	Path  string
	Alias string
}

func (n ImportNode) StmtType() string { return "ImportNode" }
