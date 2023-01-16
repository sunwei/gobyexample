package main

import (
	"bytes"
	"fmt"
	"github.com/sunwei/gobyexample/modules/template/escaper"
	"github.com/sunwei/gobyexample/modules/template/executor"
	"github.com/sunwei/gobyexample/modules/template/parser"
	"html/template"
)

func main() {
	d, err := parser.Parse("example",
		"<p><!-- HTML comment -->abc</p>\n{{.Content}}")
	if err != nil {
		fmt.Println(err)
		return
	}

	d, err = escaper.Escape(d)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := &tmpl{
		name: "hello world template",
		tmpl: d,
	}

	buf := &bytes.Buffer{}
	err = executor.Execute(t, buf, &content{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(buf.String())
}

type content struct {
}

func (c *content) Content() template.HTML {
	return template.HTML("hello template")
}

type tmpl struct {
	name string
	tmpl *parser.Document
}

func (t *tmpl) Name() string {
	return t.name
}
func (t *tmpl) Tree() *parser.Document {
	return t.tmpl
}
