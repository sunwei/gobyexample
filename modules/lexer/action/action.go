package action

import (
	"fmt"
	"github.com/sunwei/gobyexample/modules/fsm"
	"github.com/sunwei/gobyexample/modules/lexer"
)

const (
	tokenEOF lexer.TokenType = iota
	tokenLeftDelim
	tokenRightDelim
	tokenText
	tokenField
)

type token struct {
	typ lexer.TokenType
	val string
}

func (t *token) Type() lexer.TokenType {
	return t.typ
}
func (t *token) Value() string {
	return t.val
}

type delim string

const (
	left  delim = "{{"
	right delim = "}}"
)

type lex struct {
	input string
	left  delim
	right delim
	token chan *token
	fsm   fsm.FSM
}

func New(input string) lexer.Lexer {
	f := fsm.New(textState, &data{
		err: nil,
		raw: input,
	})

	l := &lex{
		input: input,
		left:  left,
		right: right,
		token: make(chan *token),
		fsm:   f,
	}

	initFSM(l)
	go l.run()

	return l
}

func (l *lex) Next() lexer.Token {
	return <-l.token
}

func (l *lex) run() {
	for {
		e := l.fsm.Process("continue")
		if e != nil {
			fmt.Println("break because of error")
			break
		}
	}
	close(l.token)
}

func (l *lex) emit(t *token) {
	l.token <- t
}

type data struct {
	err error
	raw string
}

func (d *data) Error() error {
	return d.err
}
func (d *data) Raw() any {
	return d.raw
}
