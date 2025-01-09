package lexer

type Token struct {
	Kind   TokenType
	Lexeme string
	Line   int
	Column int
}

type TokenType int

const (
	TOKEN_IDENTIFIER = iota

	TOKEN_VAR
	TOKEN_IMPORT
	TOKEN_LOOP
	TOKEN_FN
	TOKEN_IF
	TOKEN_NOTHING
	TOKEN_AS
	TOKEN_ELSE
	TOKEN_ELSEIF
	TOKEN_STRUCT
	TOKEN_FOR
	TOKEN_CONTINUE
	TOKEN_BREAK
	TOKEN_IN
	TOKEN_OR
	TOKEN_AND
	TOKEN_NOT
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_RETURN
	TOKEN_TRY
	TOKEN_CATCH
	TOKEN_FINALLY

	TOKEN_NUMBER
	TOKEN_STRING
	TOKEN_BOOLEAN
	TOKEN_LBRACKET
	TOKEN_RBRACKET
	TOKEN_RBRACE
	TOKEN_LBRACE
	TOKEN_LPAR
	TOKEN_RPAR
	TOKEN_COMMA
	TOKEN_COLON
	TOKEN_TERNARY
	TOKEN_DOT
	TOKEN_LARROW
	TOKEN_OPERATOR
	TOKEN_ASSIGN
	TOKEN_EOF
	TOKEN_EOL
	TOKEN_INIT
)

var tokenTypeLookUp = map[TokenType]string{
	TOKEN_VAR:      "var",
	TOKEN_IMPORT:   "import",
	TOKEN_LOOP:     "loop",
	TOKEN_FN:       "fn",
	TOKEN_IF:       "if",
	TOKEN_NOTHING:  "Nothing",
	TOKEN_AS:       "as",
	TOKEN_ELSE:     "else",
	TOKEN_ELSEIF:   "elseif",
	TOKEN_STRUCT:   "struct",
	TOKEN_FOR:      "for",
	TOKEN_CONTINUE: "continue",
	TOKEN_BREAK:    "break",
	TOKEN_IN:       "in",
	TOKEN_OR:       "or",
	TOKEN_AND:      "and",
	TOKEN_NOT:      "not",
	TOKEN_TRUE:     "true",
	TOKEN_FALSE:    "false",
	TOKEN_RETURN:   "return",
	TOKEN_TRY:      "try",
	TOKEN_CATCH:    "catch",
	TOKEN_FINALLY:  "finally",
	TOKEN_NUMBER:   "number",
	TOKEN_STRING:   "string",
	TOKEN_BOOLEAN:  "boolean",
	TOKEN_LBRACKET: "[",
	TOKEN_RBRACKET: "]",
	TOKEN_RBRACE:   "}",
	TOKEN_LBRACE:   "{",
	TOKEN_LPAR:     "(",
	TOKEN_RPAR:     ")",
	TOKEN_COMMA:    ",",
	TOKEN_COLON:    ":",
	TOKEN_TERNARY:  "?",
	TOKEN_DOT:      ".",
	TOKEN_LARROW:   "->",
	TOKEN_OPERATOR: "operator",
	TOKEN_ASSIGN:   "=",
	TOKEN_EOF:      "eof",
	TOKEN_EOL:      "eol",
	TOKEN_INIT:     "init",
}

func GetTokenName(tokenType TokenType) string {
	return tokenTypeLookUp[tokenType]
}

func (tokenType TokenType) String() string {
	return tokenTypeLookUp[tokenType]
}
