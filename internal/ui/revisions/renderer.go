package revisions

import (
	"fmt"
	"github.com/idursun/jjui/internal/ui/operations"
	"strings"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
)

type SegmentedRenderer struct {
	Palette       common.Palette
	IsHighlighted bool
	op            operations.Operation
}

func (s SegmentedRenderer) RenderBefore(*jj.Commit) string {
	if s.IsHighlighted && s.op.RenderPosition() == operations.RenderPositionBefore {
		return s.op.Render()
	}
	return ""
}

func (s SegmentedRenderer) RenderAfter(*jj.Commit) string {
	if s.IsHighlighted && s.op.RenderPosition() == operations.RenderPositionAfter {
		return s.op.Render()
	}
	return ""
}

func (s SegmentedRenderer) RenderGlyph(connection jj.ConnectionType, commit *jj.Commit) string {
	style := s.Palette.Normal
	opMarker := ""
	if s.IsHighlighted {
		style = s.Palette.Selected
		if s.op.RenderPosition() == operations.RenderPositionGlyph {
			opMarker = s.op.Render()
		}
	}
	return style.Render(string(connection) + opMarker)
}

func (s SegmentedRenderer) RenderTermination(connection jj.ConnectionType) string {
	return s.Palette.CommitIdRestStyle.Render(string(connection))
}

func (s SegmentedRenderer) RenderChangeId(commit *jj.Commit) string {
	hidden := ""
	if commit.Hidden {
		hidden = s.Palette.Normal.Render(" hidden")
	}

	return fmt.Sprintf("%s%s %s", s.Palette.CommitShortStyle.Render(commit.ChangeIdShort), s.Palette.CommitIdRestStyle.Render(commit.ChangeId[len(commit.ChangeIdShort):]), hidden)
}

func (s SegmentedRenderer) RenderAuthor(commit *jj.Commit) string {
	if commit.IsRoot() {
		return s.Palette.Empty.Render("root()")
	}
	return s.Palette.AuthorStyle.Render(commit.Author)
}

func (s SegmentedRenderer) RenderDate(commit *jj.Commit) string {
	return s.Palette.TimestampStyle.Render(commit.Timestamp)
}

func (s SegmentedRenderer) RenderBookmarks(commit *jj.Commit) string {
	var w strings.Builder
	if s.IsHighlighted && s.op.RenderPosition() == operations.RenderPositionBookmark {
		w.WriteString(s.op.Render())
	}
	w.WriteString(s.Palette.BookmarksStyle.Render(strings.Join(commit.Bookmarks, " ")))
	return w.String()
}

func (s SegmentedRenderer) RenderMarkers(commit *jj.Commit) string {
	if commit.Conflict {
		return s.Palette.ConflictStyle.Render("conflict")
	}
	return ""
}

func (s SegmentedRenderer) RenderDescription(commit *jj.Commit) string {
	var w strings.Builder
	if s.IsHighlighted && s.op.RenderPosition() == operations.RenderPositionDescription {
		w.WriteString(s.op.Render())
	}
	if commit.Empty {
		w.WriteString(s.Palette.Empty.Render("(empty)"))
		w.WriteString(" ")
	}
	if commit.Description == "" {
		if commit.Empty {
			w.WriteString(s.Palette.Empty.Render("(no description set)"))
		} else {
			w.WriteString(s.Palette.NonEmpty.Render("(no description set)"))
		}
	} else {
		w.WriteString(s.Palette.Normal.Render(commit.Description))
	}
	return w.String()
}
