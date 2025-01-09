package parser

type NodeType uint8
type OperatorType uint8

// INTERFACES
type Stmt interface {
	StmtType() NodeType
}

type Exp interface {
	ExpType() NodeType
}

const (
	OperatorAdd OperatorType = iota
	OperatorSubtract
	OperatorMultiply
	OperatorDivide
	OperatorModulo
	OperatorAnd
	OperatorOr
	OperatorNot
	OperatorEquals
	OperatorGreaterThan
	OperatorLessThan
	OperatorGreaterOrEqThan
	OperatorLessOrEqThan
)

const (
	NodeExpStmt NodeType = iota
	NodeNumber
	NodeString
	NodeBoolean
	NodeIdentifier
	NodeNothing
	NodeAssignment
	NodeBinaryExp
	NodeUnaryExp
	NodeCallExp
	NodeArrayExp
	NodeIndexAccessExp
	NodeDictionaryExp
	NodeObjectInitExp
	NodeMemberExp
	NodeSliceExp
	NodeTernaryExp

	NodeStructDeclaration

	NodeVarDeclaration
	NodeIfStatement
	NodeForInStatement
	NodeLoopStatement
	NodeFunctionDeclaration
	NodeAnonFunctionDeclaration
	NodeReturnStatement
	NodeBreakStatement
	NodeContinueStatement
	NodeStructMethodDeclaration
	NodeTryCatchStatement
	NodeImportStatement
	NodeBinaryComparisonExp
	NodeBinaryLogicExp
)

var NodeTypeStringLookup = map[NodeType]string{
	NodeIdentifier:              "Identifier",
	NodeMemberExp:               "Member expression",
	NodeIndexAccessExp:          "Index access expression",
	NodeExpStmt:                 "Expression statement",
	NodeNumber:                  "Number",
	NodeString:                  "String",
	NodeBoolean:                 "Boolean",
	NodeNothing:                 "Nothing",
	NodeAssignment:              "Assignment",
	NodeBinaryExp:               "Binary expression",
	NodeUnaryExp:                "Unary expression",
	NodeCallExp:                 "Call expression",
	NodeArrayExp:                "Array expression",
	NodeDictionaryExp:           "Dictionary expression",
	NodeObjectInitExp:           "Object expression",
	NodeSliceExp:                "Slice expression",
	NodeTernaryExp:              "Ternary expression",
	NodeBinaryComparisonExp:     "Binary expression",
	NodeBinaryLogicExp:          "Binary expression",
	NodeStructDeclaration:       "Declaration",
	NodeVarDeclaration:          "Declaration",
	NodeIfStatement:             "If statement",
	NodeForInStatement:          "For statement",
	NodeLoopStatement:           "Loop statement",
	NodeFunctionDeclaration:     "Function declaration",
	NodeAnonFunctionDeclaration: "Anon function declaration",
	NodeReturnStatement:         "Return statement",
	NodeBreakStatement:          "Break statement",
	NodeContinueStatement:       "Continue statement",
	NodeStructMethodDeclaration: "Struct method declaration",
	NodeTryCatchStatement:       "Try statement",
	NodeImportStatement:         "Import statement",
}

func (nt NodeType) String() string {
	return NodeTypeStringLookup[nt]
}

// EXPRESIONES
type ExpressionStmtNode struct {
	Expression Exp
}

func (e ExpressionStmtNode) StmtType() NodeType { return NodeExpStmt }

type NumberNode struct {
	Value float64
	Line  int
}

func (n NumberNode) ExpType() NodeType { return NodeNumber }

type StringNode struct {
	Value string
	Line  int
}

func (n StringNode) ExpType() NodeType { return NodeString }

type BooleanNode struct {
	Value bool
	Line  int
}

func (n BooleanNode) ExpType() NodeType { return NodeBoolean }

type IdentifierNode struct {
	Value string
	Line  int
}

func (n IdentifierNode) ExpType() NodeType { return NodeIdentifier }

type NothingNode struct {
	Line int
}

func (n NothingNode) ExpType() NodeType { return NodeNothing }

type AssignmentNode struct {
	Left     Exp
	Operator string
	Right    Exp
	Line     int
}

func (n AssignmentNode) ExpType() NodeType { return NodeAssignment }

type BinaryExpNode struct {
	Left     Exp
	Operator OperatorType
	Right    Exp
	Line     int
}

func (n BinaryExpNode) ExpType() NodeType { return NodeBinaryExp }

type BinaryComparisonExpNode struct {
	Left     Exp
	Operator OperatorType
	Right    Exp
	Line     int
}

func (n BinaryComparisonExpNode) ExpType() NodeType { return NodeBinaryComparisonExp }

type BinaryLogicExpNode struct {
	Left     Exp
	Operator OperatorType
	Right    Exp
	Line     int
}

func (n BinaryLogicExpNode) ExpType() NodeType { return NodeBinaryLogicExp }

type UnaryExpNode struct {
	Operator string
	Right    Exp
	Line     int
}

func (n UnaryExpNode) ExpType() NodeType { return NodeUnaryExp }

type CallExpNode struct {
	Args []Exp
	Name Exp
	Line int
}

func (n CallExpNode) ExpType() NodeType { return NodeCallExp }

type ArrayExpNode struct {
	Value []Exp
	Line  int
}

func (n ArrayExpNode) ExpType() NodeType { return NodeArrayExp }

type IndexAccessExpNode struct {
	Left  Exp
	Index Exp
	Line  int
}

func (n IndexAccessExpNode) ExpType() NodeType { return NodeIndexAccessExp }

type DictionaryExpNode struct {
	Value map[string]Exp
	Line  int
}

func (n DictionaryExpNode) ExpType() NodeType { return NodeDictionaryExp }

// Struct Initialization, NOT ANY OBJECT LIKE JS
type ObjectInitExpNode struct {
	Struct Exp
	Value  DictionaryExpNode
	Line   int
}

func (n ObjectInitExpNode) ExpType() NodeType { return NodeObjectInitExp }

// Struct Initialization, NOT ANY OBJECT LIKE JS
type MemberExpNode struct {
	Left   Exp
	Member string
	Line   int
}

func (n MemberExpNode) ExpType() NodeType { return NodeMemberExp }

// slice

type SliceExpNode struct {
	Left Exp
	From Exp
	To   Exp
	Line int
}

func (n SliceExpNode) ExpType() NodeType { return NodeSliceExp }

type TernaryExpNode struct {
	Condition Exp
	Left      Exp
	Right     Exp
	Line      int
}

func (n TernaryExpNode) ExpType() NodeType { return NodeTernaryExp }

// STATEMENTS

type VarDeclarationNode struct {
	Left     IdentifierNode
	Operator string
	Right    Exp
	Line     int
}

func (n VarDeclarationNode) StmtType() NodeType { return NodeVarDeclaration }

type IfStatementNode struct {
	ElseIf    []IfStatementNode
	Body      []Stmt
	ElseBody  []Stmt
	Condition Exp
	Line      int
}

func (n IfStatementNode) StmtType() NodeType { return NodeIfStatement }

type FunctionDeclarationNode struct {
	Name       string
	Body       []Stmt
	Parameters []string
	Line       int
}

func (n FunctionDeclarationNode) StmtType() NodeType { return NodeFunctionDeclaration }

type AnonFunctionDeclarationNode struct {
	Body       []Stmt
	Parameters []string
	Line       int
}

func (n AnonFunctionDeclarationNode) ExpType() NodeType { return NodeAnonFunctionDeclaration }

type StructDeclarationNode struct {
	Name       string
	Properties []string
	Line       int
}

func (n StructDeclarationNode) StmtType() NodeType { return NodeStructDeclaration }

type StructMethodDeclarationNode struct {
	Struct   string
	Function FunctionDeclarationNode
	Line     int
}

func (n StructMethodDeclarationNode) StmtType() NodeType { return NodeStructMethodDeclaration }

type ForInSatementNode struct {
	Iterator     Exp
	Body         []Stmt
	IndexVarName string
	LocalVarName string
	Line         int
}

func (n ForInSatementNode) StmtType() NodeType { return NodeForInStatement }

type BreakNode struct {
	Line int
}

func (n BreakNode) StmtType() NodeType { return NodeBreakStatement }

type ContinueNode struct {
	Line int
}

func (n ContinueNode) StmtType() NodeType { return NodeContinueStatement }

type ReturnNode struct {
	Right Exp
	Line  int
}

func (n ReturnNode) StmtType() NodeType { return NodeReturnStatement }

type TryCatchNode struct {
	Body    []Stmt
	Catch   []Stmt
	Finally []Stmt
	Line    int
}

func (n TryCatchNode) StmtType() NodeType { return NodeTryCatchStatement }

type LoopStmtNode struct {
	Body []Stmt
	Line int
}

func (n LoopStmtNode) StmtType() NodeType { return NodeLoopStatement }

type ImportNode struct {
	Line  int
	Path  string
	Alias string
}

func (n ImportNode) StmtType() NodeType { return NodeImportStatement }
