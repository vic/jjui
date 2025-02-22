package common

type SetBookmarkOperation struct{}

func (s SetBookmarkOperation) Render() string {
	return ""
}

func (s SetBookmarkOperation) RenderPosition() RenderPosition {
	return RenderPositionBookmark
}
