package describe

import (
	"jjui/internal/ui/msgs"

	"github.com/charmbracelet/bubbles/textarea"

	tea "github.com/charmbracelet/bubbletea"
)

func close() tea.Msg {
	return msgs.Close{}
}

type Model struct {
	textArea textarea.Model
}

func (m Model) Init() tea.Cmd {
	return m.textArea.Focus()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, close
		}
	}
	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.textArea.View()
}

func New() tea.Model {
	return Model{
		textArea: textarea.New(),
	}
}
