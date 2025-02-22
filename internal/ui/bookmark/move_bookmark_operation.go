package bookmark

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
)

type MoveBookmarkOperation struct {
	Overlay  tea.Model
	selected *jj.Commit
}

func (m MoveBookmarkOperation) Update(msg tea.Msg) (common.Operation, tea.Cmd) {
	var cmd tea.Cmd
	m.Overlay, cmd = m.Overlay.Update(msg)
	return MoveBookmarkOperation{Overlay: m.Overlay}, cmd
}

func (m MoveBookmarkOperation) Render() string {
	return m.Overlay.View()
}

func (m MoveBookmarkOperation) RenderPosition() common.RenderPosition {
	return common.RenderPositionAfter
}

func NewMoveBookmarkOperation(commands common.UICommands, selected *jj.Commit, width int) (common.Operation, tea.Cmd) {
	op := MoveBookmarkOperation{
		selected: selected,
		Overlay:  New(commands, selected.GetChangeId(), width),
	}
	return op, tea.Batch(commands.FetchBookmarks(selected.GetChangeId()), op.Overlay.Init())
}
