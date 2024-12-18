package parser

import (
	"evie/lexer"
)

type TokenIterator struct {
	Items []lexer.Token
	Index int
}

func (t TokenIterator) Get() lexer.Token {
	if t.IsOutOfBounds() {
		return lexer.Token{Kind: "eof"}
	}
	char := t.Items[t.Index]
	return char
}
func (t *TokenIterator) Eat() lexer.Token {
	char := t.Items[t.Index]
	t.Index++
	return char
}

func (t TokenIterator) HasNext() bool {
	return t.Index < len(t.Items)
}

func (t TokenIterator) IsOutOfBounds() bool {
	return t.Index >= len(t.Items)
}

func (t TokenIterator) GetNext() lexer.Token {
	return t.Items[t.Index+1]
}
