package bookmark

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
)

type DeleteBookmarkOperation struct {
	selected *jj.Commit
	Overlay  tea.Model
}

func (d DeleteBookmarkOperation) Init() tea.Cmd {
	return d.Overlay.Init()
}

func (d DeleteBookmarkOperation) Update(msg tea.Msg) (common.Operation, tea.Cmd) {
	var cmd tea.Cmd
	d.Overlay, cmd = d.Overlay.Update(msg)
	return DeleteBookmarkOperation{Overlay: d.Overlay}, cmd
}

func (d DeleteBookmarkOperation) RenderPosition() common.RenderPosition {
	return common.RenderPositionAfter
}

func (d DeleteBookmarkOperation) Render() string {
	return d.Overlay.View()
}

func NewDeleteBookmarkOperation(commands common.UICommands, selected *jj.Commit, width int) (common.Operation, tea.Cmd) {
	op := DeleteBookmarkOperation{
		selected: selected,
		Overlay:  NewDeleteBookmark(commands, selected.GetChangeId(), selected.Bookmarks, width),
	}
	return op, op.Overlay.Init()
}
