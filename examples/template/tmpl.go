package main

import (
	"html/template"
	"os"
)

// index html template
var indexTemplate = "<html>\n" +
	"  <body>\n" +
	"    {{.Content}}\n" +
	"  </body>\n" +
	"</html>\n"

// Post struct with exposed filed Content
type Post struct {
	Content string
}

func main() {
	// New Post with content
	// Source file could be post.md
	post := Post{"<h2>Section</h2>\n" +
		"    <p>Hello World</p>\n"}

	// New template with indexTemplate, name as "index"
	tmpl, err :=
		template.New("index").Parse(indexTemplate)
	if err != nil {
		panic(err)
	}

	// Render post with template `index`
	// write result to os.Stdout
	err = tmpl.Execute(os.Stdout, post)
	if err != nil {
		panic(err)
	}
}
