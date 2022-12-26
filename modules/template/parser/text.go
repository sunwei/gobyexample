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
	return &text{
		treeNode: &treeNode{},
		value:    token.Value(),
	}, done, nil
}

func (t *textParser) Matching(tokenType lexer.TokenType) bool {
	return tokenType == t.matchingType
}

type text struct {
	*treeNode
	value string
}

func (t *text) String() string {
	return ""
}
