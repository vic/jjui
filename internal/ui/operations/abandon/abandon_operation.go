package abandon

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	Overlay  tea.Model
	selected *jj.Commit
}

func (a Operation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	a.Overlay, cmd = a.Overlay.Update(msg)
	return Operation{Overlay: a.Overlay}, cmd
}

func (a Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionAfter
}

func (a Operation) Render() string {
	return a.Overlay.View()
}

func (a Operation) Name() string {
	return "abandon"
}

func NewOperation(context context.AppContext, selected *jj.Commit) (operations.Operation, tea.Cmd) {
	op := Operation{
		selected: selected,
		Overlay:  New(context, selected.GetChangeId()),
	}
	return op, op.Overlay.Init()
}
