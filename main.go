package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"strings"
)

type model struct {
	items  []string
	cursor int
}

func fetchLog(location string) tea.Cmd {
	return func() tea.Msg {
		// invoke command jj log in the location
		cmd := exec.Command("jj", "log")
		cmd.Dir = location
		output, err := cmd.Output()
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		// split output by new line
		lines := strings.Split(string(output), "\n")
		return logCommand(lines)
	}
}

type logCommand []string

func (m model) Init() tea.Cmd {
	return fetchLog("/Users/idursun/repositories/elixir/beach_games")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	case logCommand:
		log := []string(msg)
		m.items = log
	case tea.WindowSizeMsg:
		normal = normal.Width(msg.Width)
		selected = selected.Width(msg.Width)
	}
	return m, nil
}

var normal = lipgloss.
	NewStyle().
	Foreground(lipgloss.Color("white")).
	Background(lipgloss.Color("gray"))

var selected = lipgloss.NewStyle().
	Background(lipgloss.Color("205")).
	Inherit(normal)

func (m model) View() string {
	items := ""
	for i, item := range m.items {
		if i == m.cursor {
			items += selected.Render(item) + "\n"
		} else {
			items += normal.Render(item) + "\n"
		}
	}
	bottom := fmt.Sprintf("Cursor: %v", m.cursor)
	return items + bottom
}

func initialModel() model {
	return model{
		items:  []string{"loading", "logs"},
		cursor: 0,
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: %v", err)
		os.Exit(1)
	}
}
