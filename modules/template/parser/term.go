package parser

import (
	"github.com/sunwei/gobyexample/modules/lexer"
	"github.com/sunwei/gobyexample/modules/lexer/action"
)

type termParser struct {
}

func (p *termParser) Parse(token lexer.Token) (Node, ParseState, error) {
	switch token.Type() {
	case action.TokenField:
		f := &fieldNode{
			treeNode: &treeNode{},
			value:    token.Value(),
		}
		return f, done, nil
	default:
		panic("not supported type token yet")
	}

	return nil, done, nil
}

type fieldNode struct {
	*treeNode
	value string
}

func (t *fieldNode) String() string {
	return t.value
}
