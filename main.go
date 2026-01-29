package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)


func main() {
	// depending on what is passed then do different actions
	// for now I will just do the file picker that start a tmux session in a specific directory
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error from running model: %v\n", err)
	}
}
