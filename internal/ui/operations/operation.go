package operations

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
)

type RenderPosition int

type SetOperationMsg struct{ Operation Operation }
type OperationChangedMsg struct{ Operation Operation }

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

func SetOperation(op Operation) tea.Cmd {
	return func() tea.Msg {
		return SetOperationMsg{Operation: op}
	}
}

func OperationChanged(op Operation) tea.Cmd {
	return func() tea.Msg {
		return OperationChangedMsg{Operation: op}
	}
}
