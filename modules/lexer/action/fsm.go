package action

import (
	"errors"
	"fmt"
	"github.com/sunwei/gobyexample/modules/fsm"
	"strings"
)

const (
	textState         = "text"
	insideActionState = "action"
	leftDelimState    = "leftDelim"
	fieldState        = "field"
	rightDelimState   = "rightDelim"
	eofState          = "eof"
)

func initFSM(l *lex) {
	l.fsm.Add(textState, func(event fsm.Event) (fsm.State, fsm.Data) {
		input := event.Data().Raw().(string)

		if pos := strings.Index(input, string(left)); pos >= 0 {
			if pos > 0 {
				l.emit(&token{typ: TokenText, val: input[0:pos]})
			}
			return leftDelimState, &data{err: nil, raw: input[pos:]}
		}

		if len(input) > 0 {
			l.emit(&token{typ: TokenText, val: input})
		}

		l.emit(&token{typ: TokenEOF})
		return eofState, &data{err: nil, raw: ""}
	})

	l.fsm.Add(leftDelimState, func(event fsm.Event) (fsm.State, fsm.Data) {
		input := event.Data().Raw().(string)
		pos := len(left)
		l.emit(&token{typ: TokenLeftDelim, val: input[0:pos]})
		return insideActionState, &data{err: nil, raw: input[pos:]}
	})

	l.fsm.Add(rightDelimState, func(event fsm.Event) (fsm.State, fsm.Data) {
		input := event.Data().Raw().(string)
		pos := len(right)
		l.emit(&token{typ: TokenRightDelim, val: input[0:pos]})
		return textState, &data{err: nil, raw: input[pos:]}
	})

	l.fsm.Add(insideActionState, func(event fsm.Event) (fsm.State, fsm.Data) {
		input := event.Data().Raw().(string)
		if strings.HasPrefix(input, string(right)) {
			return rightDelimState, &data{err: nil, raw: input}
		}

		c, s := nextChar(input)
		switch c {
		case eof:
			return textState, &data{err: errors.New("unclosed action")}
		case '.':
			trimInput := input[s:]
			if len(trimInput) > 0 {
				c := trimInput[0]
				if c < '0' || '9' < c {
					return fieldState, &data{err: nil, raw: input}
				}
			}
			panic("not supported yet")
		default:
			return textState, &data{err: fmt.Errorf("unrecognized character in action: %#U", c)}
		}
	})

	l.fsm.Add(fieldState, func(event fsm.Event) (fsm.State, fsm.Data) {
		input := event.Data().Raw().(string)
		_, s := nextChar(input) // dot
		for {
			c, s2 := nextChar(input[s:])

			if !isAlphaNumeric(c) {
				break
			}
			s += s2
		}
		l.emit(&token{typ: TokenField, val: input[0:s]})
		return insideActionState, &data{err: nil, raw: input[s:]}
	})

}
