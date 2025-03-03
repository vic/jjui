package bookmark

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
	"strings"
)

var (
	applyMove = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "move bookmark"))
)

type MoveBookmarkOperation struct {
	context context.AppContext
	Overlay tea.Model
}

func (m MoveBookmarkOperation) ShortHelp() []key.Binding {
	return []key.Binding{
		cancel,
		applyMove,
	}
}

func (m MoveBookmarkOperation) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

func (m MoveBookmarkOperation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	m.Overlay, cmd = m.Overlay.Update(msg)
	return MoveBookmarkOperation{Overlay: m.Overlay}, cmd
}

func (m MoveBookmarkOperation) Render() string {
	return m.Overlay.View()
}

func (m MoveBookmarkOperation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionAfter
}

func (m MoveBookmarkOperation) Name() string {
	return "bookmark"
}

func NewMoveBookmarkOperation(context context.AppContext, selected *jj.Commit) (operations.Operation, tea.Cmd) {
	op := MoveBookmarkOperation{
		context: context,
		Overlay: New(context, selected.GetChangeId()),
	}
	return op, tea.Batch(op.load(selected.GetChangeId()), op.Overlay.Init())
}

func (m MoveBookmarkOperation) load(revision string) tea.Cmd {
	return func() tea.Msg {
		output, _ := m.context.RunCommandImmediate(jj.BookmarkList(revision))
		var bookmarks []string
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			if strings.Contains(line, "@") {
				continue
			}
			bookmarks = append(bookmarks, line)
		}
		return common.UpdateBookmarksMsg{
			Bookmarks: bookmarks,
			Revision:  revision,
		}
	}
}
