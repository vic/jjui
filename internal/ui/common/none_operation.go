package common

type None struct{}

func (n None) RendersAfter() bool {
	return false
}
