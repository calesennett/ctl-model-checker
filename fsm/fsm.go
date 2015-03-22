package fsm

type StateMachine struct {
	States      []State
	Transitions []Transition
}

type State struct {
	ID    int
	Label string
}

type Transition struct {
	from int
	to   int
}

func Parse(lines []string) StateMachine {
	states := []State{}
	transitions := []Transition{}
	return StateMachine{states, transitions}
}
