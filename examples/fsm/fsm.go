package main

import (
	"errors"
	"fmt"
	"github.com/sunwei/gobyexample/modules/fsm"
)

func main() {
	// initial fsm with init state and data
	f := fsm.New(firstState, &data{
		err: nil,
		raw: "first",
	})
	// add state with handler
	f.Add(firstState,
		func(event fsm.Event) (fsm.State, fsm.Data) {
			if event.Type() == fsm.Action {
				fmt.Println(event.Data().Raw())
			}
			return secondState, &data{
				err: nil,
				raw: "second",
			}
		})
	f.Add(secondState,
		func(event fsm.Event) (fsm.State, fsm.Data) {
			if event.Type() == fsm.Action {
				fmt.Println(event.Data().Raw())
			}
			return lastState, &data{
				err: errors.New("something wrong"),
				raw: "last",
			}
		})
	// error occurs
	f.Add(lastState,
		func(event fsm.Event) (fsm.State, fsm.Data) {
			if e := event.Data().Error(); e != nil {
				fmt.Println(e)
				return errorState, nil
			}
			// if there is no error
			// quite with eof state
			return eofState, &data{
				err: nil,
				raw: "",
			}
		})

	for {
		// send message to notify fsm start the processing
		e := f.Process("continue")
		// quit with error
		if e != nil {
			fmt.Println("break because of error")
			break
		}
		// quite for eof state
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
