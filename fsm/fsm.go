package fsm

import (
	"strconv"
	"strings"
)

type StateMachine struct {
	States      []State
	Transitions []Transition
}

type State struct {
	ID      int
	Label   string
	Initial bool
}

type Transition struct {
	from int
	to   int
}

func Parse(lines []string) StateMachine {
	states := []State{}
	transitions := []Transition{}

	for i, line := range lines {
		if line[0:6] == "STATES" {
			strStates := strings.Split(line, " ")[1]
			numStates, _ := strconv.Atoi(strStates)
			states = createStates(numStates)
		}
		if line[0:4] == "INIT" {
			curLine := i
			for isInt(lines[curLine+1]) {
				state, _ := strconv.Atoi(lines[curLine+1])
				states = markInitial(states, state)
				curLine++
			}
		}
	}
	return StateMachine{states, transitions}
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
	for _, state := range states {
		if state.ID == initState {
			state.Initial = true
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
