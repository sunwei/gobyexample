package executer

import (
	"github.com/sunwei/gobyexample/modules/template/parser"
)

func evalIdentifierNode(c context, n parser.Node) (context, error) {

	name := n.String()

	return context{
		state: stateCommand,
		rcv:   c.rcv,
		w:     c.w,
		last:  c.last,
	}, nil
}
