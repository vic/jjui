package ui

import (
	"jjui/internal/ui/revisions"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	revisions revisions.Model
}

func New() Model {
	return Model{
		revisions: revisions.New(),
	}
}

func (m Model) Init() tea.Cmd {
	var cmd tea.Cmd
	cmd = m.revisions.Init()
	return cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.revisions, cmd = m.revisions.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.revisions.View()
}
