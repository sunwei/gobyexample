package escaper

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/template/parser"
)

func Escape(doc parser.Document) (parser.Document, error) {
	var escErr error
	doc.Walk(func(n parser.Node, ws parser.WalkState) parser.WalkStatus {
		if ws == parser.WalkIn {
			fmt.Println("Walk in: ")
			fmt.Printf("%s\n", n.String())

			switch n.Type() {
			case parser.TextNode:
				escErr = escapeTextNode(n)
			case parser.ActionNode:
				escErr = escapeActionNode(n)
			}

			if escErr != nil {
				return parser.WalkStop
			}

		} else {
			fmt.Println("Walk out.")
		}
		return parser.WalkContinue
	})

	if escErr != nil {
		return parser.Document{}, escErr
	}

}
