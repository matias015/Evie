package lexer

import (
	"evie/utils"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func Tokenize(input string) []Token {

	characters := []rune(input)

	t := utils.RuneIterator{
		Items: characters,
		Index: 0,
	}

	var tokens []Token
	word := ""
	line := 1

	tokens = append(tokens, Token{
		Kind:   "init",
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
				if t.Get() == '\r' && t.GetNext() == '\n' {
					line += 1
					t.Eat()
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

				if lastAdded != "eol" && lastAdded != "rbrace" && lastAdded != "lbrace" && lastAdded != "rbracket" && lastAdded != "lbracket" && lastAdded != "comma" {
					tokens = append(tokens, Token{
						Kind:   "eol",
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

			if lastAdded != "eol" && lastAdded != "rbrace" && lastAdded != "lbrace" && lastAdded != "rbracket" && lastAdded != "lbracket" && lastAdded != "comma" {
				tokens = append(tokens, Token{
					Kind:   "eol",
					Lexeme: "eol",
					Line:   line,
					Column: 0,
				})
			}
			continue
		}

		// If it is new line
		// if token == '\r' {
		// 	if t.HasNext() && t.Get() == '\n' {
		// 		t.Eat()
		// 		t.Eat()
		// 		line += 1

		// 		lastAdded := tokens[len(tokens)-1].Kind

		// 		if lastAdded != "eol" && lastAdded != "rbrace" && lastAdded != "lbrace" && lastAdded != "rbracket" && lastAdded != "lbracket" && lastAdded != "comma" {
		// 			tokens = append(tokens, Token{
		// 				Kind:   "eol",
		// 				Lexeme: "eol",
		// 				Line:   line,
		// 				Column: 0,
		// 			})
		// 		}
		// 		continue
		// 	} else {
		// 		t.Eat()
		// 		continue
		// 	}
		// }

		// If it is a letter
		if IsAlpha(token) || token == '_' {
			t.Eat()
			word += string(token)
			for {
				if t.HasNext() && (IsAlpha(t.Get()) || isNumber(t.Get()) || t.Get() == '_' || t.Get() == '$') {
					word += string(t.Eat())
				} else {
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
				Kind:   "number",
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
					}

					if t.Get() == '\r' && t.GetNext() == '\n' {
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
				Kind:   "string",
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
					Kind:   "operator",
					Lexeme: "==",
					Line:   line,
					Column: 0,
				})
				continue
			} else {
				tokens = append(tokens, Token{
					Kind:   "assign",
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
				Kind:   "lbrace",
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
				Kind:   "rbrace",
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
				Kind:   "lbracket",
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
				Kind:   "rbracket",
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
				Kind:   "comma",
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
				Kind:   "colon",
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
				Kind:   "ternaryexp",
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
				Kind:   "lpar",
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
				Kind:   "rpar",
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
				Kind:   "dot",
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
					Kind:   "operator",
					Lexeme: firstSymbol + "=",
					Line:   line,
					Column: 0,
				})
			} else {
				tokens = append(tokens, Token{
					Kind:   "operator",
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
						Kind:   "arrowleft",
						Lexeme: "->",
						Line:   line,
						Column: 0,
					})
				} else {
					tokens = append(tokens, Token{
						Kind:   "operator",
						Lexeme: "-",
						Line:   line,
						Column: 0,
					})
				}
			} else {
				tokens = append(tokens, Token{
					Kind:   "operator",
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
		Kind:   "eof",
		Lexeme: "",
		Line:   line,
		Column: 0,
	})
	return tokens[1:]
}

func TokenFromWord(w string, l int) Token {
	Kind := "identifier"
	if w == "var" {
		Kind = w
	} else if w == "fn" {
		Kind = w
	} else if w == "if" {
		Kind = w
	} else if w == "Nothing" {
		Kind = w
	} else if w == "as" {
		Kind = w
	} else if w == "else" {
		Kind = w
	} else if w == "elseif" {
		Kind = w
	} else if w == "struct" {
		Kind = w
	} else if w == "for" {
		Kind = w
	} else if w == "continue" {
		Kind = w
	} else if w == "break" {
		Kind = w
	} else if w == "in" {
		Kind = w
	} else if w == "or" {
		Kind = w
	} else if w == "not" {
		Kind = w
	} else if w == "and" {
		Kind = w
	} else if w == "false" || w == "true" {
		Kind = "boolean"
	} else if w == "nothing" {
		Kind = w
	} else if w == "return" {
		Kind = w
	} else if w == "try" {
		Kind = w
	} else if w == "catch" {
		Kind = w
	} else if w == "finally" {
		Kind = w
	} else if w == "import" {
		Kind = w
	} else if w == "loop" {
		Kind = w
	} else {
		Kind = "identifier"
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
