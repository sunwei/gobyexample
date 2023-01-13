package executer

import (
	"github.com/sunwei/gobyexample/modules/template/parser"
)

func evalFieldNode(c context, n parser.Node) (context, error) {

	ptr := c.rcv.data()
	field := n.String()
	method := ptr.MethodByName(field[1:]) // 0 is .

	return context{
		state: stateCommand,
		rcv:   c.rcv,
		w:     c.w,
		last:  c.last,
	}, nil
}
