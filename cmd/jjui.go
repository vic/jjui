package main

import (
	"fmt"
	"os"

	"jjui/internal/ui/revisions"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(revisions.New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
