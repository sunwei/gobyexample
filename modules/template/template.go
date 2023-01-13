package template

import "github.com/sunwei/gobyexample/modules/template/parser"

type Template interface {
	Name() string
	Tree() *parser.Document
}
