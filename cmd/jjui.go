package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/ui"
	"os"
)

var Version = "unknown"

func main() {
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
