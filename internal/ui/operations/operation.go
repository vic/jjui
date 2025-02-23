package operations

import tea "github.com/charmbracelet/bubbletea"

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
