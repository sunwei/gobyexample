package executer

import (
	"github.com/sunwei/gobyexample/modules/template/parser"
	"reflect"
)

func evalActionNode(c context, n parser.Node) (context, error) {
	return context{
		state: stateAction,
		rcv:   c.rcv,
		w:     c.w,
		last:  reflect.Value{},
	}, nil
}
