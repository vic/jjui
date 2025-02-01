package abandon

import (
	"jjui/internal/ui/common"
	"jjui/internal/ui/confirmation"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	confirmation confirmation.Model
	common.Commands
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.confirmation, cmd = m.confirmation.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.confirmation.View()
}

func New(commands common.Commands, revision string) tea.Model {
	model := confirmation.New("Are you sure you want to abandon this revision?")
	model.AddOption("Yes", tea.Batch(commands.Abandon(revision), common.Close))
	model.AddOption("No", common.Close)

	return Model{
		confirmation: model,
		Commands:     commands,
	}
}
