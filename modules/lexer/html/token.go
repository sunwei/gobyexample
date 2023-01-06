package html

import "github.com/sunwei/gobyexample/modules/lexer"

type Token struct {
	lexer.BaseToken
	Start lexer.Delim
	End   lexer.Delim
}
