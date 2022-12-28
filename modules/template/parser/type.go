package parser

import "github.com/sunwei/gobyexample/modules/lexer"

type ParseState int

const (
	done ParseState = 1 + iota
	open
)

type Parser interface {
	Parse(token lexer.Token) (Node, ParseState, error)
}

type Nodes []Node

type Node interface {
	String() string
	TreeNode
}

type TreeNode interface {
	AppendChild(node TreeNode)
	Children() []TreeNode
}
