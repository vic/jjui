package common

type DeleteBookmarkOperation struct{}

func (d DeleteBookmarkOperation) RendersAfter() bool {
	return false
}
