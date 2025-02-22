package common

type ShowDetailsOperation struct{}

func (s ShowDetailsOperation) RendersAfter() bool {
	return true
}
