package bookmark

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
)

type DeleteBookmarkOperation struct {
	Overlay tea.Model
}

func (d DeleteBookmarkOperation) ShortHelp() []key.Binding {
	return []key.Binding{
		cancel,
		apply,
	}
}

func (d DeleteBookmarkOperation) FullHelp() [][]key.Binding {
	return [][]key.Binding{d.ShortHelp()}
}

func (d DeleteBookmarkOperation) Init() tea.Cmd {
	return d.Overlay.Init()
}

func (d DeleteBookmarkOperation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	d.Overlay, cmd = d.Overlay.Update(msg)
	return d, cmd
}

func (d DeleteBookmarkOperation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionAfter
}

func (d DeleteBookmarkOperation) Render() string {
	return d.Overlay.View()
}

func (d DeleteBookmarkOperation) Name() string {
	return "bookmark"
}

func NewDeleteBookmarkOperation(context context.AppContext, selected *jj.Commit) (operations.Operation, tea.Cmd) {
	op := DeleteBookmarkOperation{
		Overlay: NewDeleteBookmark(context, selected.GetChangeId(), selected.Bookmarks),
	}
	return op, op.Overlay.Init()
}
