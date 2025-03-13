package abandon

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"
	"github.com/idursun/jjui/internal/ui/context"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	confirmation tea.Model
	context      context.AppContext
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

func New(context context.AppContext, selectedRevisions []string) tea.Model {
	message := "Are you sure you want to abandon this revision?"
	if len(selectedRevisions) > 1 {
		message = fmt.Sprintf("Are you sure you want to abandon %d revisions?", len(selectedRevisions))
	}
	model := confirmation.New(message)
	model.AddOption("Yes", context.RunCommand(jj.Abandon(selectedRevisions...), common.Refresh, common.Close), key.NewBinding(key.WithKeys("y")))
	model.AddOption("No", common.Close, key.NewBinding(key.WithKeys("n", "esc")))

	return Model{
		confirmation: &model,
		context:      context,
	}
}
