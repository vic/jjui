package bookmark

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type DeleteBookmarkOperation struct {
	Overlay tea.Model
}

func (d DeleteBookmarkOperation) Init() tea.Cmd {
	return d.Overlay.Init()
}

func (d DeleteBookmarkOperation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	d.Overlay, cmd = d.Overlay.Update(msg)
	return DeleteBookmarkOperation{Overlay: d.Overlay}, cmd
}

func (d DeleteBookmarkOperation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionAfter
}

func (d DeleteBookmarkOperation) Render() string {
	return d.Overlay.View()
}

func NewDeleteBookmarkOperation(commands common.UICommands, selected *jj.Commit, width int) (operations.Operation, tea.Cmd) {
	op := DeleteBookmarkOperation{
		Overlay: NewDeleteBookmark(commands, selected.GetChangeId(), selected.Bookmarks, width),
	}
	return op, op.Overlay.Init()
}
