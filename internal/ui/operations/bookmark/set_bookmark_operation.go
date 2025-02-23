package bookmark

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type SetBookmarkOperation struct {
	Overlay tea.Model
}

func (s SetBookmarkOperation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	s.Overlay, cmd = s.Overlay.Update(msg)
	return SetBookmarkOperation{Overlay: s.Overlay}, cmd
}

func (s SetBookmarkOperation) Render() string {
	return s.Overlay.View()
}

func (s SetBookmarkOperation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionBookmark
}

func NewSetBookmarkOperation(commands common.UICommands, selected *jj.Commit) (operations.Operation, tea.Cmd) {
	op := SetBookmarkOperation{
		Overlay: NewSetBookmark(commands, selected.GetChangeId()),
	}
	return op, op.Overlay.Init()
}
