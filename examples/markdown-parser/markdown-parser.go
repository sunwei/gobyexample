package main

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/markdown/parser"
)

func main() {
	var r = newReader()
	d, _ := parser.Parse(r.lines)
	d.Walk(func(v any,
		ws parser.WalkState) parser.WalkStatus {
		b := v.(parser.Block)
		if ws == parser.WalkIn {
			fmt.Println("Walk in: ")
			fmt.Printf("%s\n", b)
		} else {
			fmt.Println("Walk out.")
		}
		return parser.WalkContinue
	})
}

type line struct {
	raw       string
	startChar string
}

func (l *line) Raw() string {
	return l.raw
}

func (l *line) StartChar() string {
	return l.startChar
}

type reader struct {
	lines []parser.Line
}

func newReader() *reader {
	return &reader{lines: []parser.Line{
		&line{raw: "### first blog\n", startChar: "#"},
		&line{raw: "Hello Blog\n", startChar: "H"},
	}}
}
