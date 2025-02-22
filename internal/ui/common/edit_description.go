package common

type EditDescriptionOperation struct{}

func (e EditDescriptionOperation) RendersAfter() bool {
	return true
}
