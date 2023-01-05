package parser

import (
	"errors"
	"github.com/sunwei/gobyexample/modules/lexer"
	"github.com/sunwei/gobyexample/modules/lexer/action"
)

func init() {
	p := &textParser{matchingType: action.TokenText}
	registerRootParsers(p.matchingType, p)
}

type textParser struct {
	matchingType lexer.TokenType
}

func (t *textParser) Parse(token lexer.Token) (Node, ParseState, error) {
	if token.Type() != t.matchingType {
		return nil, done, errors.New("mismatch token type")
	}
	return &textNode{
		treeNode: &treeNode{},
		value:    token.Value(),
	}, done, nil
}

type textNode struct {
	*treeNode
	value string
}

func (t *textNode) String() string {
	return t.value
}

func (t *textNode) Type() NodeType {
	return TextNode
}
