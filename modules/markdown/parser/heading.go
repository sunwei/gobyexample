package parser

import "strings"

type head struct {
	BaseBlock
	level int
}

type headLine struct {
	raw string
}

func (l *headLine) Raw() string {
	return l.raw
}

func (l *headLine) StartChar() string {
	panic("not implemented")
}

type heading struct {
}

func newHeading() Parser {
	return &heading{}
}

func (h *heading) Identifiers() []Identifier {
	return []Identifier{"#"}
}

func (h *heading) Open(line Line) (Block, error) {
	return &head{
		BaseBlock: BaseBlock{
			s:  Opening,
			p:  h,
			ls: []Line{&headLine{raw: headerRaw(line.Raw())}},
		},
		level: countLevel(line.Raw()),
	}, nil
}

func (h *heading) Continue(line Line, block Block) ParseState {
	return Close
}

func headerRaw(s string) string {
	return strings.TrimSpace(strings.TrimLeft(s, "#"))
}

func countLevel(s string) int {
	return len(s) - len(strings.TrimLeft(s, "#"))
}
