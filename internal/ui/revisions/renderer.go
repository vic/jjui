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
	Overlay tea.Model
)

func If(cond bool, segments ...interface{}) interface{} {
	return ifSegment{Cond: cond, Segments: segments}
}

func Separate(sep string, segments ...interface{}) interface{} {
	return separate{sep: sep, segments: segments}
}

type SegmentedRenderer struct {
	Palette             common.Palette
	HighlightedRevision string
	Overlay             tea.Model
}

func (s *SegmentedRenderer) RenderCommit(commit *jj.Commit) string {
    highlighted := commit.ChangeIdShort == s.HighlightedRevision
	return segmentedRenderer(commit, s.Palette, highlighted,
		Separate(" ", ChangeId{}, Author{}, Timestamp{}, Bookmarks{}, ConflictMarker{}), "\n",
		Separate(" ", If(s.Overlay == nil || !highlighted, If(commit.Empty, Empty{}, " "), Description{}), If(s.Overlay != nil && highlighted, s.Overlay)),
		"\n")
}

func (s *SegmentedRenderer) RenderElidedRevisions() string {
	return s.Palette.CommitIdRestStyle.Render("~  (elided revisions)")
}

func (s *SegmentedRenderer) RenderGlyph(commit *jj.Commit) string {
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	style := s.Palette.Normal
	if highlighted {
	  style = s.Palette.Selected
	}

	if commit.Immutable {
		return style.Render("◆  ")
	} else if commit.IsWorkingCopy {
		return style.Render("@  ")
	} else {
		return style.Render("○  ")
	}
}

func segmentedRenderer(commit *jj.Commit, palette common.Palette, highlighted bool, segments ...interface{}) string {
	var w strings.Builder
	segmentedCommitRenderer(&w, commit, palette, highlighted, segments...)
	return w.String()
}

func segmentedCommitRenderer(w *strings.Builder, commit *jj.Commit, palette common.Palette, highlighted bool, segments ...interface{}) {
	for _, segment := range segments {
		switch segment := segment.(type) {
		case Overlay:
			w.WriteString(segment.View())
		case separate:
			for i, s := range segment.segments {
				previousLength := w.Len()
				segmentedCommitRenderer(w, commit, palette, highlighted, s)
				written := w.Len() > previousLength
				if written && i < len(segment.segments)-1 {
					w.WriteString(segment.sep)
				}
			}
		case ifSegment:
			if segment.Cond {
				segmentedCommitRenderer(w, commit, palette, highlighted, segment.Segments...)
			}
		case ChangeId:
			w.WriteString(palette.CommitShortStyle.Render(commit.ChangeIdShort))
			w.WriteString(palette.CommitIdRestStyle.Render(commit.ChangeId[len(commit.ChangeIdShort):]))
		case Author:
			w.WriteString(palette.AuthorStyle.Render(commit.Author))
		case Timestamp:
			w.WriteString(palette.TimestampStyle.Render(commit.Timestamp))
		case Bookmarks:
			w.WriteString(palette.BookmarksStyle.Render(strings.Join(commit.Bookmarks, " ")))
		case ConflictMarker:
			if commit.Conflict {
				w.WriteString(palette.ConflictStyle.Render("conflict"))
			}
		case Empty:
			if commit.Empty {
				w.WriteString(palette.Empty.Render("(empty)"))
			}
		case Description:
			if commit.Description == "" {
				if commit.Empty {
					w.WriteString(palette.Empty.Render("(no description set)"))
				} else {
					w.WriteString(palette.NonEmpty.Render("(no description set)"))
				}
			} else {
				w.WriteString(palette.Normal.Render(commit.Description))
			}
		case string:
			w.WriteString(segment)
		}
	}
}
