package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/ui"
	"os"
)

func main() {
	p := tea.NewProgram(revisions.New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
