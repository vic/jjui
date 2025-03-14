package abandon

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	model   tea.Model
	context context.AppContext
}

func (a Operation) Update(msg tea.Msg) (operations.OperationWithOverlay, tea.Cmd) {
	var cmd tea.Cmd
	a.model, cmd = a.model.Update(msg)
	return a, cmd
}

func (a Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionAfter
}

func (a Operation) Render() string {
	return a.model.View()
}

func (a Operation) Name() string {
	return "abandon"
}

func NewOperation(context context.AppContext, selectedRevisions []string) operations.Operation {
	message := "Are you sure you want to abandon this revision?"
	if len(selectedRevisions) > 1 {
		message = fmt.Sprintf("Are you sure you want to abandon %d revisions?", len(selectedRevisions))
	}
	model := confirmation.New(message)
	model.AddOption("Yes", context.RunCommand(jj.Abandon(selectedRevisions...), common.Refresh, common.Close), key.NewBinding(key.WithKeys("y")))
	model.AddOption("No", common.Close, key.NewBinding(key.WithKeys("n", "esc")))

	op := Operation{
		model: &model,
	}
	return op
}
