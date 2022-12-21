package main

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/lexer"
	"github.com/sunwei/gobyexample/modules/lexer/action"
)

func main() {
	lex, err := action.New("<p><!-- HTML comment -->abc</p>\n{{.Content}}")
	if err != nil {
		fmt.Println(err)
		return
	}

	var tokens []lexer.Token
	for {
		token := lex.Next()
		tokens = append(tokens, token)
		if token.Type() == action.TokenEOF {
			break
		}
	}

	for i, t := range tokens {
		fmt.Println(i + 1)
		fmt.Println(t.Value())
	}

	return
}
