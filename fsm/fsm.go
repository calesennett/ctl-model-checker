package fsm

import (
	mtx "github.com/skelterjohn/go.matrix"
	"strconv"
	"strings"
)

// A StateMachine models a finite state machine with
// a set of states and transitions.
type StateMachine struct {
	States      []State
	Transitions []Transition
}

// A State can be labeled with a proposition
// and optionally be an initial state.
type State struct {
	ID      int
	Label   string
	Initial bool
}

// A transition has a from state
// and a to state.
type Transition struct {
	from int
	to   int
}

// Parse parses a .fsm file (std-in) and creates
// and returns a StateMachine and a list of
// computations to be performed on the machine.
func Parse(lines []string) (StateMachine, []string) {
	states := []State{}
	transitions := []Transition{}
	computations := []string{}
	var numStates int

	for i, line := range lines {
		if len(line) > 5 && line[0:6] == "STATES" {
			strStates := strings.Split(line, " ")[1]
			numStates, _ = strconv.Atoi(strStates)
			states = createStates(numStates)
		} else if len(line) > 3 && line[0:4] == "INIT" {
			curLine := i + 1
			for isInt(lines[curLine]) {
				state, _ := strconv.Atoi(lines[curLine])
				states = markInitial(states, state)
				curLine++
			}
		} else if len(line) > 3 && line[0:4] == "ARCS" {
			curLine := i + 1
			for curLine < len(lines) && lines[curLine] != "" {
				from, _ := strconv.Atoi(strings.Split(lines[curLine], ":")[0])
				to, _ := strconv.Atoi(strings.Split(lines[curLine], ":")[1])
				if from < numStates && to < numStates {
					transitions = append(transitions, Transition{from, to})
				}
				curLine++
			}
		} else if len(line) > 4 && line[0:5] == "LABEL" {
			curLine := i
			label := strings.Split(lines[curLine], " ")[1]
			for curLine+1 < len(lines) && isInt(lines[curLine+1]) {
				state, _ := strconv.Atoi(lines[curLine+1])
				states = labelState(states, label, state)
				curLine++
			}
		} else if len(line) > 9 && line[0:10] == "PROPERTIES" {
			curLine := i + 1
			for curLine < len(lines) {
				computations = append(computations, lines[curLine])
				curLine++
			}
		}
	}
	return StateMachine{states, transitions}, computations
}

func (sm *StateMachine) ToMatrix() *mtx.SparseMatrix {
	elems := make(map[int]float64)
	matrix := mtx.MakeSparseMatrix(elems, len(sm.States), len(sm.States))
	for _, t := range sm.Transitions {
		matrix.Set(t.to, t.from, 1)
	}
	return matrix
}

func (s *State) HasLabel(label string) bool {
	labels := strings.Split(s.Label, ",")
	for _, l := range labels {
		if l == label {
			return l == label
		}
	}
	return false
}

func createStates(numStates int) []State {
	states := []State{}
	for i := 0; i < numStates; i++ {
		state := State{i, "", false}
		states = append(states, state)
	}
	return states
}

func markInitial(states []State, initState int) []State {
	for i, state := range states {
		if state.ID == initState {
			states[i].Initial = true
		}
	}
	return states
}

func labelState(states []State, label string, state int) []State {
	for i, s := range states {
		if s.ID == state {
			if states[i].Label == "" {
				states[i].Label = label
			} else {
				states[i].Label += "," + label
			}

		}
	}
	return states
}

func isInt(state string) bool {
	if _, err := strconv.Atoi(state); err == nil {
		return true
	}
	return false
}

func (fsm StateMachine) Satisfies(res *mtx.SparseMatrix) int {
	for i, state := range fsm.States {
		if state.Initial {
			if !(res.Get(0, i) == 1) {
				return 0
			}
		}
	}
	return 1
}
