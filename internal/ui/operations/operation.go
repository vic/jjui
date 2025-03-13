package operations

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
)

type SetOperationMsg struct{ Operation Operation }
type OperationChangedMsg struct{ Operation Operation }

type RenderPosition int

const (
	RenderPositionNil RenderPosition = iota
	RenderPositionAfter
	RenderPositionBefore
	RenderPositionGlyph
	RenderPositionBookmark
	RenderPositionDescription
	RenderPositionTop
)

type Operation interface {
	RenderPosition() RenderPosition
	Render() string
	Name() string
}

type OperationWithOverlay interface {
	Operation
	Update(msg tea.Msg) (OperationWithOverlay, tea.Cmd)
}

type TracksSelectedRevision interface {
	SetSelectedRevision(commit *jj.Commit)
}

type HandleKey interface {
	HandleKey(msg tea.KeyMsg) tea.Cmd
}
