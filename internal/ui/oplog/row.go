package oplog

import (
	"github.com/idursun/jjui/internal/screen"
)

type row struct {
	OperationId string
	Lines       []*rowLine
}

type rowLine struct {
	Segments []*screen.Segment
}

func (l *rowLine) FindIdIndex() int {
	for i, segment := range l.Segments {
		if len(segment.Text) == 12 {
			return i
		}
	}
	return -1
}

func newRowLine(segments []*screen.Segment) rowLine {
	return rowLine{Segments: segments}
}
