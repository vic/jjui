package common

type Sizable interface {
	Width() int
	Height() int
	SetWidth(w int)
	SetHeight(h int)
}
