package values

type ValueType uint8

const (
	StringType ValueType = iota
	NumberType
	BoolType
	NothingType
	ErrorType
	StructType
	DictionaryType
	ContinueType
	BreakType
	FunctionType
	NativeFunctionType
	ArrayType
	ReturnType
	ObjectType
	NamespaceType
	FileType
	NativeMethodType
	CustomType
)

func (v ValueType) String() string {
	return [...]string{
		"string",
		"number",
		"boolean",
		"nothing",
		"error",
		"Struct",
		"Dictionary",
		"Continue",
		"Break",
		"Function",
		"NativeFunction",
		"array",
		"Return",
		"Object",
		"Namespace",
		"File",
		"Custom",
	}[v]
}

type RuntimeValue interface {
	GetType() ValueType
	GetProp(name string) (RuntimeValue, error)
	GetString() string
	GetNumber() float64
	GetBool() bool
}
