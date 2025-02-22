package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
)

type ShowDetailsOperation struct {
	Overlay tea.Model
}

func (s ShowDetailsOperation) Update(msg tea.Msg) (common.Operation, tea.Cmd) {
	var cmd tea.Cmd
	s.Overlay, cmd = s.Overlay.Update(msg)
	return ShowDetailsOperation{Overlay: s.Overlay}, cmd
}

func (s ShowDetailsOperation) Render() string {
	return s.Overlay.View()
}

func (s ShowDetailsOperation) RenderPosition() common.RenderPosition {
	return common.RenderPositionAfter
}

func Operation(commands common.UICommands, selected *jj.Commit) (common.Operation, tea.Cmd) {
	op := ShowDetailsOperation{
		Overlay: New(selected.ChangeId, commands),
	}
	return op, op.Overlay.Init()
}
