package main

import (
	"bufio"
	algs "github.com/calesennett/ctl-model-checker/algorithms"
	fsm "github.com/calesennett/ctl-model-checker/fsm"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	stateMachine, computations := fsm.Parse(lines)
	algs.Run(stateMachine, computations)
}
