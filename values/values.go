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
	CapturedErrorType
	FileType
	CustomType
)

func (v ValueType) String() string {
	return [...]string{
		"string",
		"number",
		"boolean",
		"nothing",
		"error",
		"StructType",
		"DictionaryType",
		"ContinueType",
		"BreakType",
		"FunctionType",
		"NativeFunctionType",
		"array",
		"ReturnType",
		"ObjectType",
		"NamespaceType",
		"CapturedErrorType",
		"FileType",
		"CustomType",
	}[v]
}

type RuntimeValue interface {
	GetType() ValueType
	GetProp(v *RuntimeValue, name string) (RuntimeValue, error)
	GetString() string
	GetNumber() float64
	GetBool() bool
}
