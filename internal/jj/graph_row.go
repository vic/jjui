package jj

import "strings"

type GraphRow struct {
	Commit       *Commit
	IsSelected   bool
	IsAffected   bool
	SegmentLines []SegmentedLine
	Indent       int
	Previous     *GraphRow
}

func NewGraphRow() GraphRow {
	return GraphRow{
		Commit:       &Commit{},
		SegmentLines: make([]SegmentedLine, 0),
	}
}

func (r *GraphRow) AddLine(line SegmentedLine) {
	if r.Commit == nil {
		return
	}
	switch len(r.SegmentLines) {
	case 0:
		line.Flags = Revision | Highlightable
		if r.Commit == nil {
			break
		}
		r.Commit.IsWorkingCopy = line.ContainsRune('@', r.Indent)
		for i := line.ChangeIdIdx; i < line.CommitIdIdx; i++ {
			segment := line.Segments[i]
			if strings.TrimSpace(segment.Text) == "hidden" {
				r.Commit.Hidden = true
			}
		}
	default:
		if line.ContainsRune('~', r.Indent) {
			line.Flags = Elided
		} else {
			lastLine := r.SegmentLines[len(r.SegmentLines)-1]
			line.Flags = lastLine.Flags & ^Revision & ^Elided
		}
	}
	r.SegmentLines = append(r.SegmentLines, line)
}

func (r *GraphRow) Last(flag SegmentedLineFlag) *SegmentedLine {
	for i := len(r.SegmentLines) - 1; i >= 0; i-- {
		if r.SegmentLines[i].Flags&flag == flag {
			return &r.SegmentLines[i]
		}
	}
	return &SegmentedLine{}
}

type SegmentedLineIteratorPredicate func(f SegmentedLineFlag) bool

func Including(flags SegmentedLineFlag) SegmentedLineIteratorPredicate {
	return func(f SegmentedLineFlag) bool {
		return f&flags == flags
	}
}

func Excluding(flags SegmentedLineFlag) SegmentedLineIteratorPredicate {
	return func(f SegmentedLineFlag) bool {
		return f&flags != flags
	}
}

func (r *GraphRow) SegmentLinesIter(predicate SegmentedLineIteratorPredicate) func(yield func(line *SegmentedLine) bool) {
	return func(yield func(line *SegmentedLine) bool) {
		for i := range r.SegmentLines {
			line := &r.SegmentLines[i]
			if predicate(line.Flags) {
				if !yield(line) {
					return
				}
			}
		}
	}
}
