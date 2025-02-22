package common

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
