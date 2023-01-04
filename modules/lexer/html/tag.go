package html

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/lexer"
)

func readTagName(input string) (string, int) {
	pos := 0
	for {
		c, s := lexer.NextChar(input[pos:])
		fmt.Println("-=-=-=-", string(c))
		switch c {
		case ' ':
			panic("attribute not support yet")
		case '>':
			return input[:pos], pos
		}
		pos += s
	}
}
