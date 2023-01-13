package executer

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/template"
	"github.com/sunwei/gobyexample/modules/template/parser"
	"io"
	"reflect"
)

func Execute(t template.Template, w io.Writer, data any) error {
	var exeErr error
	c := context{
		state: stateText,
		rcv:   newReceiver(data),
		w:     w,
		last:  reflect.Value{},
	}

	doc := t.Tree()
	doc.Walk(func(n parser.Node, ws parser.WalkState) parser.WalkStatus {
		if ws == parser.WalkIn {
			switch n.Type() {
			case parser.TextNode:
				c.state = stateText
			case parser.ActionNode:
				c, exeErr = evalActionNode(c, n)
			case parser.CommandNode:
				c, exeErr = evalCommandNode(c, n)
			case parser.FieldNode:
				c, exeErr = evalFieldNode(c, n)
			case parser.IdentifierNode:
				c, exeErr = evalIdentifierNode(c, n)
			}

			if exeErr != nil {
				return parser.WalkStop
			}
		} else if ws == parser.WalkOut {
			switch n.Type() {
			case parser.TextNode:
				if _, err := c.w.Write([]byte(n.String())); err != nil {
					panic(fmt.Sprintf("%s: text node write error %#v", t.Name(), err))
				}
			case parser.ActionNode:

			}
		}
		return parser.WalkContinue
	})

	if exeErr != nil {
		return exeErr
	}

	return nil
}
