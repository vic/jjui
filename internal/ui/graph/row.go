package graph

import (
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/screen"
	"strings"
)

type Row struct {
	Commit     *jj.Commit
	Lines      []*GraphRowLine
	IsSelected bool
	IsAffected bool
	Indent     int
	Previous   *Row
}

type SegmentedLineFlag int

const (
	Revision SegmentedLineFlag = 1 << iota
	Highlightable
	Elided
)

type GraphRowLine struct {
	Segments    []*screen.Segment
	Flags       SegmentedLineFlag
	ChangeIdIdx int
	CommitIdIdx int
}

func NewGraphRowLine(segments []*screen.Segment) GraphRowLine {
	return GraphRowLine{
		Segments:    segments,
		ChangeIdIdx: -1,
		CommitIdIdx: -1,
	}
}

func (gr *GraphRowLine) Extend(indent int) GraphRowLine {
	ret := NewGraphRowLine(make([]*screen.Segment, 0))
	for _, s := range gr.Segments {
		extended := screen.Segment{
			Params: s.Params,
		}
		text := ""
		for _, p := range s.Text {
			if p == '│' || p == '╭' || p == '├' || p == '┐' || p == '┤' || p == '┌' { // curved, square style
				p = '│'
			} else if p == '|' { //ascii style
				p = '|'
			} else {
				p = ' '
			}
			indent--
			text += string(p)
			if indent <= 0 {
				break
			}
		}
		extended.Text = text
		ret.Segments = append(ret.Segments, &extended)
		if indent <= 0 {
			break
		}
	}
	for indent > 0 {
		ret.Segments[len(ret.Segments)-1].Text += " "
		indent--
	}
	return ret
}
func (gr *GraphRowLine) FindIdIndex(start int) int {
	for i := start; i < len(gr.Segments); i++ {
		cur := gr.Segments[i].Text
		if cur != "" && !strings.Contains(cur, " ") {
			n := i + 1
			if n < len(gr.Segments) {
				cur = gr.Segments[n].Text
				if cur != "" && !strings.Contains(cur, " ") {
					return i
				}
			}
		}
	}
	return -1
}

func (gr *GraphRowLine) ContainsRune(r rune, indent int) bool {
	for _, segment := range gr.Segments {
		text := segment.Text
		if len(segment.Text) > indent {
			text = segment.Text[:indent]
		}
		indent -= len(text)
		if strings.ContainsRune(text, r) {
			return true
		}
	}
	return false
}

func NewGraphRow() Row {
	return Row{
		Commit: &jj.Commit{},
		Lines:  make([]*GraphRowLine, 0),
	}
}

func (r *Row) AddLine(line *GraphRowLine) {
	if r.Commit == nil {
		return
	}
	switch len(r.Lines) {
	case 0:
		line.Flags = Revision | Highlightable
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
			lastLine := r.Lines[len(r.Lines)-1]
			line.Flags = lastLine.Flags & ^Revision & ^Elided
		}
	}
	r.Lines = append(r.Lines, line)
}

func (r *Row) Last(flag SegmentedLineFlag) *GraphRowLine {
	for i := len(r.Lines) - 1; i >= 0; i-- {
		if r.Lines[i].Flags&flag == flag {
			return r.Lines[i]
		}
	}
	return &GraphRowLine{}
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

func (r *Row) SegmentLinesIter(predicate SegmentedLineIteratorPredicate) func(yield func(line *GraphRowLine) bool) {
	return func(yield func(line *GraphRowLine) bool) {
		for i := range r.Lines {
			line := r.Lines[i]
			if predicate(line.Flags) {
				if !yield(line) {
					return
				}
			}
		}
	}
}
