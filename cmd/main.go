package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"jjui/internal/jj"
	"os"
	"strings"
)

type model struct {
	items    []jj.Commit
	selected map[int]bool
	cursor   int
	width    int
}

func fetchLog(location string) tea.Cmd {
	return func() tea.Msg {
		lines := jj.GetCommits(location)
		return logCommand(lines)
	}
}

type logCommand []jj.Commit

func (m model) Init() tea.Cmd {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return fetchLog(dir)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "m", " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		default:
			return m, nil
		}
	case logCommand:
		commits := []jj.Commit(msg)
		m.items = commits
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

var commitShortStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ff00ff"))

var commitShortStyleHighlighted = lipgloss.NewStyle().
	Background(lipgloss.Color("#425893")).
	Inherit(commitShortStyle)

var commitIdRestStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#c0c0c0"))

var commitIdRestHighlightedStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#425893")).
	Inherit(commitIdRestStyle)

var normal = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ffffff"))

var normalHighlighted = lipgloss.NewStyle().
	Background(lipgloss.Color("#425893")).
	Inherit(normal)

func (m model) View() string {
	items := ""
	for i, commit := range m.items {
		changeIdRemaining := strings.TrimPrefix(commit.ChangeId, commit.ChangeIdShort)
		item := ""
		if r, ok := m.selected[i]; ok && r {
			items += normalHighlighted.Render("x")
		} else {
			items += normal.Render(".")
		}
		if i == m.cursor {
			item += commitShortStyleHighlighted.Render(commit.ChangeIdShort)
			item += commitIdRestHighlightedStyle.Render(changeIdRemaining + " ")
			item += normalHighlighted.Width(m.width).Render(commit.Description) + "\n"
		} else {
			item += commitShortStyle.Render(commit.ChangeIdShort)
			item += commitIdRestStyle.Render(changeIdRemaining + " ")
			item += normal.Render(commit.Description) + "\n"
		}
		items += item
	}
	bottom := fmt.Sprintf("use j,k keys to move and down: %v", m.cursor)

	return items + bottom
}

func initialModel() model {
	return model{
		items:    []jj.Commit{},
		selected: make(map[int]bool),
		cursor:   0,
		width:    20,
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: %v", err)
		os.Exit(1)
	}
}
