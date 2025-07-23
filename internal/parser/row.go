package parser

import (
	"strings"
	"unicode"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/screen"
)

type Row struct {
	Commit     *jj.Commit
	Lines      []*GraphRowLine
	IsAffected bool
	Indent     int
	Previous   *Row
}

type RowLineFlags int

const (
	Revision RowLineFlags = 1 << iota
	Highlightable
	Elided
)

type GraphGutter struct {
	Segments []*screen.Segment
}

type GraphRowLine struct {
	Segments []*screen.Segment
	Gutter   GraphGutter
	Flags    RowLineFlags
}

func NewGraphRowLine(segments []*screen.Segment) GraphRowLine {
	return GraphRowLine{
		Segments: segments,
		Gutter:   GraphGutter{Segments: make([]*screen.Segment, 0)},
	}
}

func (gr *GraphRowLine) chop(indent int) {
	if len(gr.Segments) == 0 {
		return
	}
	segments := gr.Segments
	gr.Segments = make([]*screen.Segment, 0)

	for i, s := range segments {
		extended := screen.Segment{
			Style: s.Style,
		}
		var textBuilder strings.Builder
		for _, p := range s.Text {
			if indent <= 0 {
				break
			}
			textBuilder.WriteRune(p)
			indent--
		}
		extended.Text = textBuilder.String()
		gr.Gutter.Segments = append(gr.Gutter.Segments, &extended)
		if len(extended.Text) < len(s.Text) {
			gr.Segments = append(gr.Segments, &screen.Segment{
				Text:  s.Text[len(extended.Text):],
				Style: s.Style,
			})
		}
		if indent <= 0 && len(segments)-i-1 > 0 {
			gr.Segments = segments[i+1:]
			break
		}
	}

	// Pad with spaces if indent is not fully consumed
	if indent > 0 && len(gr.Segments) > 0 {
		lastSegment := gr.Segments[len(gr.Segments)-1]
		lastSegment.Text += strings.Repeat(" ", indent)
	}

	// break gutter into segments per rune
	segments = gr.Gutter.Segments
	gr.Gutter.Segments = make([]*screen.Segment, 0)
	for _, s := range segments {
		for _, p := range s.Text {
			extended := screen.Segment{
				Text:  string(p),
				Style: s.Style,
			}
			gr.Gutter.Segments = append(gr.Gutter.Segments, &extended)
		}
	}
}

func (gr *GraphRowLine) containsRune(r rune) bool {
	for _, segment := range gr.Gutter.Segments {
		if strings.ContainsRune(segment.Text, r) {
			return true
		}
	}
	return false
}

func isChangeIdLike(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func isHexLike(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.Hex_Digit, r) {
			return false
		}
	}
	return true
}

func (gr *GraphRowLine) FindPossibleChangeIdIdx() int {
	for i, segment := range gr.Segments {
		if isChangeIdLike(segment.Text) {
			return i
		}
	}
	return -1
}

func (gr *GraphRowLine) FindPossibleCommitIdIdx(after int) int {
	for i := after; i < len(gr.Segments); i++ {
		segment := gr.Segments[i]
		if isHexLike(segment.Text) {
			return i
		}
	}
	return -1
}

func NewGraphRow() Row {
	return Row{
		Commit: &jj.Commit{},
		Lines:  make([]*GraphRowLine, 0),
	}
}

func (row *Row) Extend() GraphGutter {
	type extendResult int
	const (
		No extendResult = iota
		Yes
		Carry
	)
	canExtend := func(text string) extendResult {
		for _, p := range text {
			switch p {
			case '│', '|', '╭', '├', '┐', '┤', '┌', '╮', '┬', '┼', '+', '\\', '.':
				return Yes
			case '╯', '╰', '└', '┴', '┘', ' ', '/':
				return No
			case '─', '-':
				return Carry
			}
		}
		return No
	}

	extendMask := make([]bool, len(row.Lines[0].Gutter.Segments))
	var lastGutter *GraphGutter
	for _, gl := range row.Lines {
		if gl.Flags&Highlightable != Highlightable {
			continue
		}
		for i, s := range gl.Gutter.Segments {
			answer := canExtend(s.Text)
			switch answer {
			case Yes:
				extendMask[i] = true
			case No:
				extendMask[i] = false
			case Carry:
				extendMask[i] = extendMask[i]
			}
		}
		lastGutter = &gl.Gutter
	}

	if lastGutter == nil {
		return GraphGutter{Segments: make([]*screen.Segment, 0)}
	}

	if len(extendMask) > len(lastGutter.Segments) {
		extendMask = extendMask[:len(lastGutter.Segments)]
	}
	ret := GraphGutter{
		Segments: make([]*screen.Segment, len(extendMask)),
	}
	for i, b := range extendMask {
		ret.Segments[i] = &screen.Segment{
			Style: lastGutter.Segments[i].Style,
			Text:  " ",
		}
		if b {
			ret.Segments[i].Text = "│"
		}
	}
	return ret
}

func (row *Row) AddLine(line *GraphRowLine) {
	if row.Commit == nil {
		return
	}
	line.chop(row.Indent)
	switch len(row.Lines) {
	case 0:
		line.Flags = Revision | Highlightable
		row.Commit.IsWorkingCopy = line.containsRune('@')
		for _, segment := range line.Segments {
			if strings.TrimSpace(segment.Text) == "hidden" {
				row.Commit.Hidden = true
			}
		}
	default:
		if line.containsRune('~') {
			line.Flags = Elided
		} else {
			if row.Commit.CommitId == "" {
				commitIdIdx := line.FindPossibleCommitIdIdx(0)
				if commitIdIdx != -1 {
					row.Commit.CommitId = line.Segments[commitIdIdx].Text
					line.Flags = Revision | Highlightable
				}
			}
			lastLine := row.Lines[len(row.Lines)-1]
			line.Flags = lastLine.Flags & ^Revision & ^Elided
		}
	}
	row.Lines = append(row.Lines, line)
}

func (row *Row) Last(flag RowLineFlags) *GraphRowLine {
	for i := len(row.Lines) - 1; i >= 0; i-- {
		if row.Lines[i].Flags&flag == flag {
			return row.Lines[i]
		}
	}
	return &GraphRowLine{}
}

func (row *Row) RowLinesIter(predicate RowLinesIteratorPredicate) func(yield func(line *GraphRowLine) bool) {
	return func(yield func(line *GraphRowLine) bool) {
		for i := range row.Lines {
			line := row.Lines[i]
			if predicate(line.Flags) {
				if !yield(line) {
					return
				}
			}
		}
	}
}

type RowLinesIteratorPredicate func(f RowLineFlags) bool

func Including(flags RowLineFlags) RowLinesIteratorPredicate {
	return func(f RowLineFlags) bool {
		return f&flags == flags
	}
}

func Excluding(flags RowLineFlags) RowLinesIteratorPredicate {
	return func(f RowLineFlags) bool {
		return f&flags != flags
	}
}
