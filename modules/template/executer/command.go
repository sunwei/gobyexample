package executer

import (
	"github.com/sunwei/gobyexample/modules/template/parser"
)

func evalCommandNode(c context, n parser.Node) (context, error) {
	return context{
		state: stateCommand,
		rcv:   c.rcv,
		w:     c.w,
		last:  c.last,
	}, nil
}
