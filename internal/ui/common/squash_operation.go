package common

type SquashOperation struct {
	From string
}

func (s SquashOperation) RendersAfter() bool {
	return false
}
