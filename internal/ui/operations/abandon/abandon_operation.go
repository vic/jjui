package abandon

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	Overlay  tea.Model
	selected []*jj.Commit
}

func (a Operation) Update(msg tea.Msg) (operations.OperationWithOverlay, tea.Cmd) {
	var cmd tea.Cmd
	a.Overlay, cmd = a.Overlay.Update(msg)
	return a, cmd
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

func NewOperation(context context.AppContext, selectedRevisions []*jj.Commit) (operations.Operation, tea.Cmd) {
	var changeIds []string
	for _, s := range selectedRevisions {
		changeIds = append(changeIds, s.GetChangeId())
	}
	op := Operation{
		selected: selectedRevisions,
		Overlay:  New(context, changeIds),
	}
	return op, op.Overlay.Init()
}
