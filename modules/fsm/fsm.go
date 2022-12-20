package fsm

type Data interface {
	Error() error
	Raw() any
}
type State string
type StateHandler func(event Event) (State, Data)

type TypeEvent int

const (
	Notification TypeEvent = iota
	Action
)

type Event interface {
	Type() TypeEvent
	Message() any
	Data() Data
}

type FSM interface {
	Add(state State, handler StateHandler)
	Process(message any) error
	State() State
}

func New(initState State, initData Data) FSM {
	return &fsm{
		state:    initState,
		data:     initData,
		handlers: map[State]StateHandler{},
	}
}

type fsm struct {
	state    State
	data     Data
	handlers map[State]StateHandler
}

func (f *fsm) State() State {
	return f.state
}

func (f *fsm) Add(state State, handler StateHandler) {
	if _, ok := f.handlers[state]; ok {
		panic("state handler exist already")
	}
	f.handlers[state] = handler
}

func (f *fsm) Process(message any) error {
	h, ok := f.handlers[f.state]
	if !ok {
		panic("state handler not exist")
	}
	s, d := h(&event{
		t:       Action,
		data:    f.data,
		message: message,
	})
	f.state = s
	f.data = d
	return f.data.Error()
}

type event struct {
	t       TypeEvent
	data    Data
	message any
}

func (e *event) Type() TypeEvent {
	return e.t
}
func (e *event) Data() Data {
	return e.data
}
func (e *event) Message() any {
	return e.message
}
func (e *event) Error() error {
	return e.data.Error()
}
