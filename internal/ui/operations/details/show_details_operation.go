package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	Overlay tea.Model
}

func (s Operation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	s.Overlay, cmd = s.Overlay.Update(msg)
	return Operation{Overlay: s.Overlay}, cmd
}

func (s Operation) Render() string {
	return s.Overlay.View()
}

func (s Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionAfter
}

func NewOperation(commands common.UICommands, selected *jj.Commit) (operations.Operation, tea.Cmd) {
	op := Operation{
		Overlay: New(selected.ChangeId, commands),
	}
	return op, op.Overlay.Init()
}
