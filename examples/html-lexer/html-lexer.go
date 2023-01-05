package main

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/lexer"
	"github.com/sunwei/gobyexample/modules/lexer/html"
)

func main() {
	lex, err := html.New("<p><!-- HTML comment -->abc</p>")
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
		if token.Type() == html.TokenEOF {
			break
		}
	}

	// output tokens detail
	for i, t := range tokens {
		fmt.Println(i + 1)
		fmt.Println(t.Value())
	}
}
