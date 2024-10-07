package status

import (
	"jjui/internal/ui/common"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	spinner spinner.Model
	command string
	output  string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.CommandRunningMsg:
		m.command = string(msg)
		return m, m.spinner.Tick
	case common.CommandCompletedMsg:
		m.command = ""
		m.output = msg.Output
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	s := m.spinner.View()
	if m.command == "" {
		s = ""
	}
	return s + " " + m.command
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{
		spinner: s,
		command: "",
		output:  "",
	}
}
