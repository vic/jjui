package revisions

import (
	"jjui/internal/jj"
	"strings"

	"jjui/internal/ui/common"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	ChangeId       struct{}
	Author         struct{}
	Timestamp      struct{}
	Bookmarks      struct{}
	ConflictMarker struct{}
	Empty          struct{}
	Description    struct{}
	separate       struct {
		sep      string
		segments []interface{}
	}
	ifSegment struct {
		Cond     bool
		Segments []interface{}
	}
	NodeGlyph       struct{}
	Glyph           struct{}
	ElidedRevisions struct{}
	Overlay         tea.Model
)

func If(cond bool, segments ...interface{}) interface{} {
	return ifSegment{Cond: cond, Segments: segments}
}

func Separate(sep string, segments ...interface{}) interface{} {
	return separate{sep: sep, segments: segments}
}

func SegmentedRenderer(w *strings.Builder, row *jj.GraphRow, palette common.Palette, highlighted bool, segments ...interface{}) {
	for _, segment := range segments {
		switch segment := segment.(type) {
		case Overlay:
			w.WriteString(segment.View())
		case separate:
			for i, s := range segment.segments {
				previousLength := w.Len()
				SegmentedRenderer(w, row, palette, highlighted, s)
				written := w.Len() > previousLength
				if written && i < len(segment.segments)-1 {
					w.WriteString(segment.sep)
				}
			}
		case ifSegment:
			if segment.Cond {
				SegmentedRenderer(w, row, palette, highlighted, segment.Segments...)
			}
		case ElidedRevisions:
			if row.Elided {
				indent := strings.Repeat("│ ", row.Level)
				w.WriteString(indent)
				w.WriteString(palette.CommitIdRestStyle.Render("~ (elided revisions)"))
				w.WriteString("\n")
			}
		case NodeGlyph:
			nodeGlyph := "○"
			switch {
			case row.Commit.IsWorkingCopy:
				nodeGlyph = "@"
			case row.Commit.Immutable:
				nodeGlyph = "◆"
			case row.Commit.Conflict:
				nodeGlyph = "×"
			}
			indent := strings.Repeat("│ ", row.Level)
			w.WriteString(indent)
			if highlighted {
				w.WriteString(palette.Selected.Render(nodeGlyph))
			} else {
				w.WriteString(nodeGlyph)
			}
		case Glyph:
			indent := strings.Repeat("│ ", row.Level)
			glyph := "│"
			if len(row.Node.Parents) > 0 && len(row.Node.Parents[0].Edges) > 1 && row.Level > 0 {
				glyph = "├─╯"
				indent = strings.Repeat("│ ", row.Level-1)
			}
			w.WriteString(indent)
			w.WriteString(glyph)
		case ChangeId:
			w.WriteString(palette.CommitShortStyle.Render(row.Commit.ChangeIdShort))
			w.WriteString(palette.CommitIdRestStyle.Render(row.Commit.ChangeId[len(row.Commit.ChangeIdShort):]))
		case Author:
			w.WriteString(palette.AuthorStyle.Render(row.Commit.Author))
		case Timestamp:
			w.WriteString(palette.TimestampStyle.Render(row.Commit.Timestamp))
		case Bookmarks:
			w.WriteString(palette.BookmarksStyle.Render(strings.Join(row.Commit.Bookmarks, " ")))
		case ConflictMarker:
			if row.Commit.Conflict {
				w.WriteString(palette.ConflictStyle.Render("conflict"))
			}
		case Empty:
			if row.Commit.Empty {
				w.WriteString(palette.Empty.Render("(empty)"))
			}
		case Description:
			if row.Commit.Description == "" {
				if row.Commit.Empty {
					w.WriteString(palette.Empty.Render("(no description set)"))
				} else {
					w.WriteString(palette.NonEmpty.Render("(no description set)"))
				}
			} else {
				w.WriteString(palette.Normal.Render(row.Commit.Description))
			}
		case string:
			w.WriteString(segment)
		}
	}
}
