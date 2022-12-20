package action

import (
	"errors"
	"github.com/sunwei/gobyexample/modules/fsm"
	"strings"
)

const (
	textState       = "text"
	leftDelimState  = "leftDelim"
	fieldState      = "field"
	rightDelimState = "rightDelim"
)

var inputTypeErr = errors.New("input not string type")
var inputEmptyErr = errors.New("input empty")

func initFSM(l *lex) {
	l.fsm.Add(textState, func(event fsm.Event) (fsm.State, fsm.Data) {
		input, ok := event.Data().Raw().(string)
		if !ok {
			return textState, &data{
				err: inputTypeErr,
			}
		}

		if pos := strings.Index(input, string(left)); pos >= 0 {
			if pos > 0 {
				l.emit(&token{typ: tokenText, val: input[0:pos]})
			}
			return leftDelimState, &data{
				err: nil,
				raw: input[pos:],
			}
		}
		if len(input) > 0 {
			l.emit(&token{typ: tokenText, val: input})
		}

		l.emit(&token{typ: tokenEOF})
		return textState, &data{
			err: inputEmptyErr,
		}
	})

	l.fsm.Add(leftDelimState, func(event fsm.Event) (fsm.State, fsm.Data) {
		input, ok := event.Data().Raw().(string)
		if !ok {
			return textState, &data{
				err: inputTypeErr,
			}
		}

	})

}
