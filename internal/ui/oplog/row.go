package oplog

import "github.com/idursun/jjui/internal/screen"

type Row struct {
	OperationId string
	Lines       []*RowLine
}

type RowLine struct {
	Segments []*screen.Segment
}

func (l *RowLine) FindIdIndex() int {
	for i, segment := range l.Segments {
		if len(segment.Text) == 12 {
			return i
		}
	}
	return -1
}

func NewRowLine(segments []*screen.Segment) RowLine {
	return RowLine{Segments: segments}
}
