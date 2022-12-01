package parser

type Line interface {
	Raw() string
	StartChar() string
}

type Identifier string

type ParseState int

const (
	Continue ParseState = 1 << iota
	Children
	Close
)

type Parser interface {
	Identifiers() []Identifier
	Open(line Line) (Block, error)
	Continue(line Line, block Block) ParseState
}

type Market interface {
	FindParser(identifier Identifier) Parser
}

type BlockState int

const (
	Opening BlockState = 1 << iota
	Closed
)

type Block interface {
	IsOpen() bool
	Close()
	Parser() Parser
	AppendLine(line Line)
}

type BaseBlock struct {
	s  BlockState
	p  Parser
	ls []Line
}

func (b *BaseBlock) IsOpen() bool {
	return b.s == Opening
}

func (b *BaseBlock) Close() {
	b.s = Closed
}

func (b *BaseBlock) Parser() Parser {
	return b.p
}

func (b *BaseBlock) AppendLine(line Line) {
	b.ls = append(b.ls, line)
}

type Inline struct {
	p  Parser
	ls []Line
}

type WalkStatus int

const (
	WalkStop WalkStatus = iota + 1
	WalkContinue
)

type WalkState int

const (
	WalkIn WalkState = 1 << iota
	WalkOut
)

type Walker func(v any, ws WalkState) WalkStatus
