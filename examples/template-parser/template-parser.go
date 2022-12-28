package main

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/template/parser"
)

func main() {
	d, err := parser.Parse("example",
		"<p><!-- HTML comment -->abc</p>\n{{.Content}}")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d.String())
}
