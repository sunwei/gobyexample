package parser

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/lexer"
	"github.com/sunwei/gobyexample/modules/lexer/action"
)

func Parse(name string, text string) (*tree, error) {
	lex, err := action.New(text)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	p := &parser{
		name:  name,
		lexer: lex,
		tree:  newTree(),
	}

	err = p.parse()
	if err != nil {
		return nil, err
	}

	return p.tree, err
}

var rootParsers map[lexer.TokenType]Parser

func registerRootParsers(tokenType lexer.TokenType, p Parser) {
	if _, ok := rootParsers[tokenType]; ok {
		panic("duplicated parser")
	}
	rootParsers[tokenType] = p
}

func getParser(tokenType lexer.TokenType) Parser {
	return rootParsers[tokenType]
}

type parser struct {
	name  string
	lexer lexer.Lexer
	tree  *tree
}

func (p *parser) parse() error {
	var currentParser Parser
	ps := done

	for {
		token := p.lexer.Next()
		if token.Type() == action.TokenEOF {
			break
		}

		// keep the same parser only after it's done
		if ps == done {
			currentParser = getParser(token.Type())
		}

		n, ps2, err := currentParser.Parse(token)
		if err != nil {
			return err
		}
		ps = ps2
		p.tree.root.AppendChild(n)
	}

	return nil
}
