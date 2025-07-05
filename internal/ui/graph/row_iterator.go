package graph

import (
	"io"
)

type RowIterator interface {
	Len() int
	Next() bool
	Render(w io.Writer)
	RowHeight() int
	IsHighlighted() bool
}
