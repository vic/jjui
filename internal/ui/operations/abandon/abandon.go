package abandon

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	confirmation tea.Model
	context      common.AppContext
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

func New(context common.AppContext, revision string) tea.Model {
	model := confirmation.New("Are you sure you want to abandon this revision?")
	model.AddOption("Yes", tea.Batch(context.RunCommand(jj.Abandon(revision), common.Refresh("@")), common.Close), key.NewBinding(key.WithKeys("y")))
	model.AddOption("No", common.Close, key.NewBinding(key.WithKeys("n", "esc")))

	return Model{
		confirmation: &model,
		context:      context,
	}
}
