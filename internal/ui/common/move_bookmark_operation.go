package common

type MoveBookmarkOperation struct{}

func (m MoveBookmarkOperation) RendersAfter() bool {
	return false
}
