package escaper

import (
	"bytes"
	"github.com/sunwei/gobyexample/modules/lexer/html"
	"github.com/sunwei/gobyexample/modules/template/parser"
	"io"
	"strings"
)

func escapeTextNode(n parser.Node) error {
	lex, err := html.New(n.String())
	if err != nil {
		return err
	}

	escapedStr := ""
	for {
		token := lex.Next()
		if token.Type() == html.TokenEOF {
			break
		}

		switch token.Type() {
		case html.TokenComment:
			// escape comment
			continue
		default:
			escapedStr += string(token.(*html.Token).Start) +
				escapeString(token.Value()) +
				string(token.(*html.Token).End)
		}
	}

	n.SetVal(escapedStr)

	return nil

}

const escapedChars = "&'<>\"\r"

func escapeString(s string) string {
	if strings.IndexAny(s, escapedChars) == -1 {
		return s
	}
	var buf bytes.Buffer
	err := escape(&buf, s)
	if err != nil {
		panic("escape error")
	}
	return buf.String()
}

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

func escape(w writer, s string) error {
	i := strings.IndexAny(s, escapedChars)
	for i != -1 {
		if _, err := w.WriteString(s[:i]); err != nil {
			return err
		}
		var esc string
		switch s[i] {
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := w.WriteString(esc); err != nil {
			return err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := w.WriteString(s)
	return err
}
