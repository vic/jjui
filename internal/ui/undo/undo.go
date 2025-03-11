package undo

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"
	"github.com/idursun/jjui/internal/ui/context"
)

type Model struct {
	confirmation tea.Model
}

func (m Model) Init() tea.Cmd {
	return m.confirmation.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.confirmation, cmd = m.confirmation.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.confirmation.View()
}

func NewModel(context context.AppContext) Model {
	model := confirmation.New("Are you sure you want to undo last change?")
	model.SetBorderStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(2))
	model.AddOption("Yes", context.RunCommand(jj.Undo(), common.Refresh, common.Close), key.NewBinding(key.WithKeys("y")))
	model.AddOption("No", common.Close, key.NewBinding(key.WithKeys("n", "esc")))
	return Model{
		confirmation: &model,
	}
}
