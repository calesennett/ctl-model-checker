package main

import (
	"bufio"
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
	stateMachine := fsm.Parse(lines)
}
