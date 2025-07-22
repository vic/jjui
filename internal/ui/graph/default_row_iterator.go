package graph

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/parser"
	"github.com/idursun/jjui/internal/ui/common"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/operations"
)

type DefaultRowIterator struct {
	SearchText    string
	Selections    map[string]bool
	Op            operations.Operation
	Width         int
	Rows          []parser.Row
	isHighlighted bool
	isSelected    bool
	current       int
	Cursor        int
	dimmedStyle   lipgloss.Style
	checkStyle    lipgloss.Style
	textStyle     lipgloss.Style
	selectedStyle lipgloss.Style
}

type Option func(*DefaultRowIterator)

func NewDefaultRowIterator(rows []parser.Row, options ...Option) *DefaultRowIterator {
	iterator := &DefaultRowIterator{
		Op:         &operations.Default{},
		Rows:       rows,
		Selections: make(map[string]bool),
		current:    -1,
	}

	for _, opt := range options {
		opt(iterator)
	}

	return iterator
}

func WithWidth(width int) Option {
	return func(s *DefaultRowIterator) {
		s.Width = width
	}
}

func WithStylePrefix(prefix string) Option {
	if prefix != "" {
		prefix = " "
	}
	return func(s *DefaultRowIterator) {
		s.textStyle = common.DefaultPalette.Get(prefix + "text").Inline(true)
		s.selectedStyle = common.DefaultPalette.Get(prefix + "selected").Inline(true)
		s.dimmedStyle = common.DefaultPalette.Get(prefix + "dimmed")
		s.checkStyle = common.DefaultPalette.Get(prefix + "success").Inline(true)
	}
}

func (s *DefaultRowIterator) IsHighlighted() bool {
	return s.current == s.Cursor
}

func (s *DefaultRowIterator) Next() bool {
	s.current++
	if s.current >= len(s.Rows) {
		return false
	}
	s.isHighlighted = s.current == s.Cursor
	s.isSelected = false
	if v, ok := s.Selections[s.Rows[s.current].Commit.GetChangeId()]; ok {
		s.isSelected = v
	}
	return true
}

func (s *DefaultRowIterator) RowHeight() int {
	return len(s.Rows[s.current].Lines)
}

func (s *DefaultRowIterator) Render(r io.Writer) {
	row := s.Rows[s.current]
	// will render by extending the previous connections
	if before := s.RenderBefore(row.Commit); before != "" {
		extended := parser.GraphRowLine{}
		if row.Previous != nil {
			extended = row.Previous.Last(parser.Highlightable).Extend(row.Indent)
		}
		s.writeSection(r, extended, row.Indent, false, before)
	}
	var lastLine *parser.GraphRowLine
	for segmentedLine := range row.RowLinesIter(parser.Including(parser.Highlightable)) {
		lastLine = segmentedLine
		lw := strings.Builder{}
		if segmentedLine.Flags&parser.Revision != parser.Revision && s.isHighlighted {
			if decoration := s.Op.Render(row.Commit, operations.RenderOverDescription); decoration != "" {
				extended := segmentedLine.Chop(row.Indent)
				s.writeSection(r, extended, row.Indent, true, decoration)
				continue
			}
		}
		// if it is a revision line
		for i, segment := range segmentedLine.Segments {
			if i == segmentedLine.ChangeIdIdx {
				if decoration := s.RenderBeforeChangeId(row.Commit); decoration != "" {
					fmt.Fprint(&lw, decoration)
				}
			}

			if s.isHighlighted && i == segmentedLine.CommitIdIdx {
				if decoration := s.RenderBeforeCommitId(row.Commit); decoration != "" {
					fmt.Fprint(&lw, decoration)
				}
			}

			style := segment.Style
			if s.isHighlighted {
				style = style.Inherit(s.selectedStyle)
			} else {
				style = style.Inherit(s.textStyle)
			}

			start, end := segment.FindSubstringRange(s.SearchText)
			if start != -1 {
				mid := lipgloss.NewRange(start, end, style.Reverse(true))
				fmt.Fprint(&lw, lipgloss.StyleRanges(style.Render(segment.Text), mid))
			} else {
				fmt.Fprint(&lw, style.Render(segment.Text))
			}
		}
		if segmentedLine.Flags&parser.Revision == parser.Revision && row.IsAffected {
			style := s.dimmedStyle
			if s.isHighlighted {
				style = s.dimmedStyle.Background(s.selectedStyle.GetBackground())
			}
			fmt.Fprint(&lw, style.Render(" (affected by last operation)"))
		}
		line := lw.String()
		if s.isHighlighted && segmentedLine.Flags&parser.Highlightable == parser.Highlightable {
			fmt.Fprint(r, lipgloss.PlaceHorizontal(s.Width, 0, line, lipgloss.WithWhitespaceBackground(s.selectedStyle.GetBackground())))
		} else {
			fmt.Fprint(r, lipgloss.PlaceHorizontal(s.Width, 0, line, lipgloss.WithWhitespaceBackground(s.textStyle.GetBackground())))
		}
		fmt.Fprint(r, "\n")
	}

	if row.Commit.IsRoot() {
		return
	}

	afterSection := s.RenderAfter(row.Commit)
	if afterSection != "" && lastLine != nil {
		s.writeSection(r, lastLine.Extend(row.Indent), row.Indent, false, afterSection)
	}

	for segmentedLine := range row.RowLinesIter(parser.Excluding(parser.Highlightable)) {
		var lw strings.Builder
		for _, segment := range segmentedLine.Segments {
			fmt.Fprint(&lw, segment.Style.Inherit(s.textStyle).Render(segment.Text))
		}
		line := lw.String()
		fmt.Fprint(r, lipgloss.PlaceHorizontal(s.Width, 0, line, lipgloss.WithWhitespaceBackground(s.textStyle.GetBackground())))
		fmt.Fprint(r, "\n")
	}
}

func (s *DefaultRowIterator) writeSection(r io.Writer, extended parser.GraphRowLine, indent int, highlight bool, section string) {
	lines := strings.Split(section, "\n")
	for _, sectionLine := range lines {
		lw := strings.Builder{}
		for _, segment := range extended.Segments {
			if s.isHighlighted && highlight {
				fmt.Fprint(&lw, segment.Style.Inherit(s.selectedStyle).Render(segment.Text))
			} else {
				fmt.Fprint(&lw, segment.Style.Inherit(s.textStyle).Render(segment.Text))
			}
		}

		fmt.Fprint(&lw, sectionLine)
		line := lw.String()
		if s.isHighlighted && highlight {
			fmt.Fprint(r, lipgloss.PlaceHorizontal(s.Width, 0, line, lipgloss.WithWhitespaceBackground(s.selectedStyle.GetBackground())))
		} else {
			fmt.Fprint(r, lipgloss.PlaceHorizontal(s.Width, 0, line, lipgloss.WithWhitespaceBackground(s.textStyle.GetBackground())))
		}
		fmt.Fprintln(r)
		extended = extended.Extend(indent)
	}
}

func (s *DefaultRowIterator) Len() int {
	return len(s.Rows)
}

func (s *DefaultRowIterator) RenderBefore(commit *jj.Commit) string {
	return s.Op.Render(commit, operations.RenderPositionBefore)
}

func (s *DefaultRowIterator) RenderAfter(commit *jj.Commit) string {
	return s.Op.Render(commit, operations.RenderPositionAfter)
}

func (s *DefaultRowIterator) RenderBeforeChangeId(commit *jj.Commit) string {
	opMarker := s.Op.Render(commit, operations.RenderBeforeChangeId)
	selectedMarker := ""
	if s.isSelected {
		if s.isHighlighted {
			selectedMarker = s.checkStyle.Background(s.selectedStyle.GetBackground()).Render("✓")
		} else {
			selectedMarker = s.checkStyle.Background(s.textStyle.GetBackground()).Render("✓")
		}
	}
	if opMarker == "" && selectedMarker == "" {
		return ""
	}
	var sections []string

	space := s.textStyle.Render(" ")
	if s.isHighlighted {
		space = s.selectedStyle.Render(" ")
	}

	if opMarker != "" {
		sections = append(sections, opMarker, space)
	}
	if selectedMarker != "" {
		sections = append(sections, selectedMarker, space)
	}
	return lipgloss.JoinHorizontal(0, sections...)
}

func (s *DefaultRowIterator) RenderBeforeCommitId(commit *jj.Commit) string {
	return s.Op.Render(commit, operations.RenderBeforeCommitId)
}
