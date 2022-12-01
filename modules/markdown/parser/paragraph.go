package parser

import "strings"

type paragraph struct {
}

func newParagraph() Parser {
	return &paragraph{}
}

func (p *paragraph) Identifiers() []Identifier {
	return nil
}

func (p *paragraph) Open(line Line) (Block, error) {
	return &BaseBlock{
		s:  Opening,
		p:  p,
		ls: []Line{line},
	}, nil
}

func (p *paragraph) Continue(line Line, block Block) ParseState {
	if isBlank(line.Raw()) {
		return Close
	}
	block.AppendLine(line)
	return Continue
}

func isBlank(s string) bool {
	tsl := strings.TrimSpace(s)
	if len(tsl) == 0 {
		return true
	}
	return false
}
