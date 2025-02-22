package common

type MoveBookmarkOperation struct{}

func (m MoveBookmarkOperation) Render() string {
	return ""
}

func (m MoveBookmarkOperation) RenderPosition() RenderPosition {
	return RenderPositionAfter
}
