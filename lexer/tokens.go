package lexer

type Token struct {
	Kind   string
	Lexeme string
	Line   int
	Column int
}
