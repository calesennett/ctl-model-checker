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

func Parse(lines []string) (StateMachine, []string) {
	states := []State{}
	transitions := []Transition{}
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
			for curLine < len(lines) && isInt(string(lines[curLine][0])) {
				from, _ := strconv.Atoi(strings.Split(lines[curLine], ":")[0])
				to, _ := strconv.Atoi(strings.Split(lines[curLine], ":")[1])
				if from < numStates && to < numStates {
					transitions = append(transitions, Transition{from, to})
				}
				curLine++
			}
		}
	}
	return StateMachine{states, transitions}, []string{}
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

func isInt(state string) bool {
	if _, err := strconv.Atoi(state); err == nil {
		return true
	}
	return false
}
