package abandon

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
)

type AbandonOperation struct {
	Overlay  tea.Model
	selected *jj.Commit
}

func (a AbandonOperation) Update(msg tea.Msg) (common.Operation, tea.Cmd) {
	var cmd tea.Cmd
	a.Overlay, cmd = a.Overlay.Update(msg)
	return AbandonOperation{Overlay: a.Overlay}, cmd
}

func (a AbandonOperation) RenderPosition() common.RenderPosition {
	return common.RenderPositionAfter
}

func (a AbandonOperation) Render() string {
	return a.Overlay.View()
}

func Operation(commands common.UICommands, selected *jj.Commit) (common.Operation, tea.Cmd) {
	op := AbandonOperation{
		selected: selected,
		Overlay:  New(commands, selected.GetChangeId()),
	}
	return op, op.Overlay.Init()
}
