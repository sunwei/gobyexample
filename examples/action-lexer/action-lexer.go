package main

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/lexer"
	"github.com/sunwei/gobyexample/modules/lexer/action"
)

func main() {
	// Action example
	lex, err := action.New(
		"<p><!-- HTML comment -->abc</p>\n{{.Content}}")
	if err != nil {
		fmt.Println(err)
		return
	}

	var tokens []lexer.Token
	for {
		// lexer iterate
		token := lex.Next()
		tokens = append(tokens, token)
		// reach end, analyzing done
		if token.Type() == action.TokenEOF {
			break
		}
	}

	// output tokens detail
	for i, t := range tokens {
		fmt.Println(i + 1)
		fmt.Println(t.Value())
	}

	return
}
