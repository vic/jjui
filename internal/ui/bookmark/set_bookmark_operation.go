package bookmark

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
)

type SetBookmarkOperation struct {
	selected *jj.Commit
	Overlay  tea.Model
}

func (s SetBookmarkOperation) Update(msg tea.Msg) (common.Operation, tea.Cmd) {
	var cmd tea.Cmd
	s.Overlay, cmd = s.Overlay.Update(msg)
	return SetBookmarkOperation{Overlay: s.Overlay}, cmd
}

func (s SetBookmarkOperation) Render() string {
	return s.Overlay.View()
}

func (s SetBookmarkOperation) RenderPosition() common.RenderPosition {
	return common.RenderPositionBookmark
}

func NewSetBookmarkOperation(commands common.UICommands, selected *jj.Commit) (common.Operation, tea.Cmd) {
	op := SetBookmarkOperation{
		selected: selected,
		Overlay:  NewSetBookmark(commands, selected.GetChangeId()),
	}
	return op, op.Overlay.Init()
}
