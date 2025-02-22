package common

type Operation interface {
	RendersAfter() bool
}
