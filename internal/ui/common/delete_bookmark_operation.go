package common

type DeleteBookmarkOperation struct{}

func (d DeleteBookmarkOperation) RenderPosition() RenderPosition {
	return RenderPositionAfter
}

func (d DeleteBookmarkOperation) Render() string {
	return ""
}
