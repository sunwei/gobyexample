package parser

import (
	"github.com/sunwei/gobyexample/modules/lexer"
)

func newCommand(tokens lexer.Tokens) (*commandNode, error) {
	tp := &termParser{}
	cmd := &commandNode{treeNode: &treeNode{}}
	for _, t := range tokens {
		n, _, err := tp.Parse(t)
		if err != nil {
			return nil, err
		}
		cmd.AppendChild(n)
	}

	return cmd, nil
}

type commandNode struct {
	*treeNode
}

func (n *commandNode) String() string {
	cs := n.Children()
	s := ""
	for _, n := range cs {
		s += n.(Node).String()
		s += " "
	}
	return s
}

func (n *commandNode) Type() NodeType {
	return CommandNode
}
