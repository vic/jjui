package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/ui"
	"os"
)

var Version = "unknown"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("jjui version %s\n", Version)
		os.Exit(0)
	}
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
