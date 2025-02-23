package operations

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
)

type RenderPosition int

const (
	RenderPositionNil RenderPosition = iota
	RenderPositionAfter
	RenderPositionBefore
	RenderPositionGlyph
	RenderPositionBookmark
	RenderPositionDescription
)

type Operation interface {
	RenderPosition() RenderPosition
	Render() string
}

type OperationWithOverlay interface {
	Update(msg tea.Msg) (Operation, tea.Cmd)
}

type TracksSelectedRevision interface {
	SetSelectedRevision(commit *jj.Commit)
}

type HandleKey interface {
	HandleKey(msg tea.KeyMsg) tea.Cmd
}
