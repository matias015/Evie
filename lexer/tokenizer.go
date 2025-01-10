package lexer

import (
	"evie/common"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func Tokenize(input string) []Token {

	characters := []rune(input)

	t := common.RuneIterator{
		Items: characters,
		Index: 0,
	}

	var tokens []Token
	word := ""
	line := 1

	tokens = append(tokens, Token{
		Kind:   TOKEN_INIT,
		Lexeme: "init",
		Line:   line,
		Column: 0,
	})

	for {
		if t.IsOutOfBounds() {
			break
		}

		token := t.Get()

		// If it is space
		if token == ' ' || token == '\t' {
			t.Eat()
			continue
		}

		// If it is comment
		if token == '/' && t.HasNext() && t.GetNext() == '/' {
			t.Eat()
			t.Eat()
			for {
				if t.Get() == '\r' {
					t.Eat()
					continue
				} else if t.Get() == '\n' {
					line += 1
					t.Eat()
					break
				} else {
					t.Eat()
				}
			}
			continue
		}

		if token == '\r' {
			if t.HasNext() && t.GetNext() == '\n' {
				t.Eat()
				t.Eat()
				line += 1

				lastAdded := tokens[len(tokens)-1].Kind

				if lastAdded != TOKEN_EOL && lastAdded != TOKEN_RBRACE && lastAdded != TOKEN_LBRACE && lastAdded != TOKEN_RBRACKET && lastAdded != TOKEN_LBRACKET && lastAdded != TOKEN_COMMA {
					tokens = append(tokens, Token{
						Kind:   TOKEN_EOL,
						Lexeme: "eol",
						Line:   line,
						Column: 0,
					})
				}
				continue
			} else {
				t.Eat()
				continue
			}
		}

		if token == '\n' {
			t.Eat()
			line += 1
			lastAdded := tokens[len(tokens)-1].Kind

			if lastAdded != TOKEN_EOL && lastAdded != TOKEN_RBRACE && lastAdded != TOKEN_LBRACE && lastAdded != TOKEN_RBRACKET && lastAdded != TOKEN_LBRACKET && lastAdded != TOKEN_COMMA {
				tokens = append(tokens, Token{
					Kind:   TOKEN_EOL,
					Lexeme: "eol",
					Line:   line,
					Column: 0,
				})
			}
			continue
		}

		// If it is a letter
		if IsAlpha(token) || token == '_' {
			word += string(token)
			for {
				if t.HasNext() && (IsAlpha(t.GetNext()) || isNumber(t.GetNext()) || t.GetNext() == '_' || t.GetNext() == '$') {
					t.Eat()
					word += string(t.Get())
				} else {
					t.Eat()
					break
				}
			}

			tokens = append(tokens, TokenFromWord(word, line))
			word = ""
			continue
		}

		// If it is a number
		if isNumber(token) {
			t.Eat()
			word += string(token)
			for {
				if t.HasNext() && (isNumber(t.Get()) || t.Get() == '.') {
					word += string(t.Eat())
				} else {
					break
				}
			}

			tokens = append(tokens, Token{
				Kind:   TOKEN_NUMBER,
				Lexeme: word,
				Line:   line,
				Column: 0,
			})
			word = ""
			continue
		}

		// If it is a string
		if token == '"' {
			t.Eat()

			initLine := line

			for {

				if t.HasNext() && t.Get() != '"' {

					if string(t.Get()) == "\\" && string(t.GetNext()) == "n" {
						t.Eat()
						t.Eat()
						word += string('\n')
						continue
					} else if string(t.Get()) == "\\" && string(t.GetNext()) == "r" {
						t.Eat()
						t.Eat()
						word += string('\r')
						continue
					} else if string(t.Get()) == "\\" && string(t.GetNext()) == "t" {
						t.Eat()
						t.Eat()
						word += string('\t')
						continue
					}

					if t.Get() == '\r' {
						t.Eat()
						continue
					}

					if t.Get() == '\n' {
						line += 1
					}
					word += string(t.Eat())
				} else {
					if t.Get() != '"' {
						fmt.Println("string started at line " + strconv.Itoa(initLine) + " not closed")
						os.Exit(1)
					}
					t.Eat()
					break
				}
			}

			tokens = append(tokens, Token{
				Kind:   TOKEN_STRING,
				Lexeme: word,
				Line:   line,
				Column: 0,
			})

			word = ""
			continue
		}

		// if is an = or ==
		if token == '=' {
			t.Eat()
			if t.HasNext() && t.Get() == '=' {
				t.Eat()
				tokens = append(tokens, Token{
					Kind:   TOKEN_OPERATOR,
					Lexeme: "==",
					Line:   line,
					Column: 0,
				})
				continue
			} else {
				tokens = append(tokens, Token{
					Kind:   TOKEN_ASSIGN,
					Lexeme: "=",
					Line:   line,
					Column: 0,
				})
				continue
			}
		}

		// }
		if token == '{' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_LBRACE,
				Lexeme: "{",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// }
		if token == '}' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_RBRACE,
				Lexeme: "}",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// [
		if token == '[' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_LBRACKET,
				Lexeme: "[",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// ]
		if token == ']' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_RBRACKET,
				Lexeme: "]",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// comma
		if token == ',' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_COMMA,
				Lexeme: ",",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// colon
		if token == ':' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_COLON,
				Lexeme: ":",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// ternaryexp
		if token == '?' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_TERNARY,
				Lexeme: "?",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// lpar
		if token == '(' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_LPAR,
				Lexeme: "(",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// rpar
		if token == ')' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_RPAR,
				Lexeme: ")",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// dot
		if token == '.' {
			t.Eat()
			tokens = append(tokens, Token{
				Kind:   TOKEN_DOT,
				Lexeme: ".",
				Line:   line,
				Column: 0,
			})
			continue
		}

		// < and > and <= and >=
		if token == '<' || token == '>' {
			firstSymbol := string(t.Eat())
			if t.HasNext() && t.Get() == '=' {
				t.Eat()
				tokens = append(tokens, Token{
					Kind:   TOKEN_OPERATOR,
					Lexeme: firstSymbol + "=",
					Line:   line,
					Column: 0,
				})
			} else {
				tokens = append(tokens, Token{
					Kind:   TOKEN_OPERATOR,
					Lexeme: firstSymbol,
					Line:   line,
					Column: 0,
				})
			}
			continue
		}

		// + - * / and also check if after an '-' there is a > to get a '->'
		if token == '+' || token == '-' || token == '*' || token == '/' {
			t.Eat()
			if token == '-' {
				if t.HasNext() && t.Get() == '>' {
					t.Eat()
					tokens = append(tokens, Token{
						Kind:   TOKEN_LARROW,
						Lexeme: "->",
						Line:   line,
						Column: 0,
					})
				} else {
					tokens = append(tokens, Token{
						Kind:   TOKEN_OPERATOR,
						Lexeme: "-",
						Line:   line,
						Column: 0,
					})
				}
			} else {
				tokens = append(tokens, Token{
					Kind:   TOKEN_OPERATOR,
					Lexeme: string(token),
					Line:   line,
					Column: 0,
				})
			}
			continue
		}

		fmt.Println("Unknown token: " + string(token) + " in line " + fmt.Sprint(line))
		os.Exit(1)

		if t.HasNext() {
			t.Eat()
		}

	}

	tokens = append(tokens, Token{
		Kind:   TOKEN_EOF,
		Lexeme: "",
		Line:   line,
		Column: 0,
	})
	return tokens[1:]
}

func TokenFromWord(w string, l int) Token {
	var Kind TokenType
	if w == "var" {
		Kind = TOKEN_VAR
	} else if w == "fn" {
		Kind = TOKEN_FN
	} else if w == "if" {
		Kind = TOKEN_IF
	} else if w == "Nothing" {
		Kind = TOKEN_NOTHING
	} else if w == "as" {
		Kind = TOKEN_AS
	} else if w == "else" {
		Kind = TOKEN_ELSE
	} else if w == "elseif" {
		Kind = TOKEN_ELSEIF
	} else if w == "struct" {
		Kind = TOKEN_STRUCT
	} else if w == "for" {
		Kind = TOKEN_FOR
	} else if w == "continue" {
		Kind = TOKEN_CONTINUE
	} else if w == "break" {
		Kind = TOKEN_BREAK
	} else if w == "in" {
		Kind = TOKEN_IN
	} else if w == "or" {
		Kind = TOKEN_OR
	} else if w == "not" {
		Kind = TOKEN_NOT
	} else if w == "and" {
		Kind = TOKEN_AND
	} else if w == "false" || w == "true" {
		Kind = TOKEN_BOOLEAN
	} else if w == "nothing" {
		Kind = TOKEN_NOTHING
	} else if w == "return" {
		Kind = TOKEN_RETURN
	} else if w == "try" {
		Kind = TOKEN_TRY
	} else if w == "catch" {
		Kind = TOKEN_CATCH
	} else if w == "finally" {
		Kind = TOKEN_FINALLY
	} else if w == "import" {
		Kind = TOKEN_IMPORT
	} else if w == "loop" {
		Kind = TOKEN_LOOP
	} else {
		Kind = TOKEN_IDENTIFIER
	}

	return Token{
		Kind:   Kind,
		Lexeme: w,
		Line:   l,
		Column: 0,
	}
}

func IsAlpha(char rune) bool {
	return unicode.IsLetter(char)
}

func isNumber(char rune) bool {
	return unicode.IsDigit(char)
}
