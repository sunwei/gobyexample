package main

import (
	"errors"
	"fmt"
	"github.com/sunwei/gobyexample/modules/fsm"
)

func main() {
	f := fsm.New(firstState, &data{
		err: nil,
		raw: "first",
	})
	f.Add(firstState, func(event fsm.Event) (fsm.State, fsm.Data) {
		if event.Type() == fsm.Action {
			fmt.Println(event.Data().Raw())
		}
		return secondState, &data{
			err: nil,
			raw: "second",
		}
	})
	f.Add(secondState, func(event fsm.Event) (fsm.State, fsm.Data) {
		if event.Type() == fsm.Action {
			fmt.Println(event.Data().Raw())
		}
		return lastState, &data{
			err: errors.New("something wrong"),
			raw: "last",
		}
	})
	f.Add(lastState, func(event fsm.Event) (fsm.State, fsm.Data) {
		if e := event.Data().Error(); e != nil {
			fmt.Println(e)
			return errorState, nil
		}
		return eofState, &data{
			err: nil,
			raw: "",
		}
	})

	for {
		e := f.Process("continue")
		if e != nil {
			fmt.Println("break because of error")
			break
		}
		if f.State() == eofState {
			fmt.Println("eof")
			break
		}
	}
}

const (
	firstState  = "first"
	secondState = "second"
	lastState   = "last"
	errorState  = "error"
	eofState    = "eof"
)

type data struct {
	err error
	raw any
}

func (d *data) Error() error {
	return d.err
}
func (d *data) Raw() any {
	return d.raw
}
