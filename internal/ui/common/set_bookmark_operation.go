package common

type SetBookmarkOperation struct{}

func (s SetBookmarkOperation) RendersAfter() bool {
	return false
}
