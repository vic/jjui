package bookmark

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
)

type SetBookmarkOperation struct {
	Overlay tea.Model
}

func (s SetBookmarkOperation) IsFocused() bool {
	return true
}

func (s SetBookmarkOperation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	s.Overlay, cmd = s.Overlay.Update(msg)
	return s, cmd
}

func (s SetBookmarkOperation) Render() string {
	return s.Overlay.View()
}

func (s SetBookmarkOperation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionBookmark
}

func (s SetBookmarkOperation) Name() string {
	return "bookmark"
}

func NewSetBookmarkOperation(context context.AppContext, selected *jj.Commit) (operations.Operation, tea.Cmd) {
	op := SetBookmarkOperation{
		Overlay: NewSetBookmark(context, selected.GetChangeId()),
	}
	return op, op.Overlay.Init()
}
