package ui

import (
	"jjui/internal/ui/describe"
	"jjui/internal/ui/msgs"
	"jjui/internal/ui/revisions"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	models []tea.Model
}

func New() Model {
	return Model{
		models: []tea.Model{revisions.New()},
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, model := range m.models {
		cmds = append(cmds, model.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case msgs.ShowDescribe:
		m.models = append(m.models, describe.New())
		return m, nil
	case msgs.Close:
		m.models = m.models[:len(m.models)-1]
		return m, nil
	}
	var cmd tea.Cmd
	top := m.models[len(m.models)-1]
	top, cmd = top.Update(msg)
	m.models[len(m.models)-1] = top
	return m, cmd
}

func (m Model) View() string {
	top := m.models[len(m.models)-1]
	return top.View()
}
