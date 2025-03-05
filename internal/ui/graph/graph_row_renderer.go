package graph

import (
	"fmt"
	"github.com/idursun/jjui/internal/ui/common"
	"strings"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/operations"
)

type SegmentedRenderer struct {
	Palette       common.Palette
	IsHighlighted bool
	Op            operations.Operation
}

func (s SegmentedRenderer) RenderBefore(*jj.Commit) string {
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionBefore {
		return s.Op.Render()
	}
	return ""
}

func (s SegmentedRenderer) RenderAfter(*jj.Commit) string {
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionAfter {
		return s.Op.Render()
	}
	return ""
}

func (s SegmentedRenderer) RenderGlyph(connection jj.ConnectionType, commit *jj.Commit) string {
	style := s.Palette.Normal
	opMarker := ""
	if s.IsHighlighted {
		style = s.Palette.Selected
		if s.Op.RenderPosition() == operations.RenderPositionGlyph {
			opMarker = s.Op.Render()
		}
	}
	return style.Render(string(connection) + opMarker)
}

func (s SegmentedRenderer) RenderTermination(connection jj.ConnectionType) string {
	return s.Palette.Elided.Render(string(connection))
}

func (s SegmentedRenderer) RenderChangeId(commit *jj.Commit) string {
	hidden := ""
	if commit.Hidden {
		hidden = s.Palette.Normal.Render(" hidden")
	}

	return fmt.Sprintf("%s%s %s", s.Palette.ChangeId.Render(commit.ChangeIdShort), s.Palette.Rest.Render(commit.ChangeId[len(commit.ChangeIdShort):]), hidden)
}

func (s SegmentedRenderer) RenderCommitId(commit *jj.Commit) string {
	if commit.IsRoot() {
		return ""
	}
	return s.Palette.CommitId.Render(commit.CommitIdShort) + s.Palette.Rest.Render(commit.CommitId[len(commit.ChangeIdShort):])
}

func (s SegmentedRenderer) RenderAuthor(commit *jj.Commit) string {
	if commit.IsRoot() {
		return s.Palette.EmptyPlaceholder.Render("root()")
	}
	return s.Palette.Author.Render(commit.Author)
}

func (s SegmentedRenderer) RenderDate(commit *jj.Commit) string {
	if commit.IsRoot() {
		return ""
	}
	return s.Palette.Timestamp.Render(commit.Timestamp)
}

func (s SegmentedRenderer) RenderBookmarks(commit *jj.Commit) string {
	var w strings.Builder
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionBookmark {
		w.WriteString(s.Op.Render())
	}
	if len(commit.Bookmarks) > 0 {
		w.WriteString(s.Palette.Bookmarks.Render(strings.Join(commit.Bookmarks, " ")))
	}
	return w.String()
}

func (s SegmentedRenderer) RenderMarkers(commit *jj.Commit) string {
	if commit.Conflict {
		return s.Palette.Conflict.Render("conflict")
	}
	return ""
}

func (s SegmentedRenderer) RenderDescription(commit *jj.Commit) string {
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionDescription {
		return s.Op.Render()
	}
	var w strings.Builder
	if commit.Empty {
		w.WriteString(s.Palette.EmptyPlaceholder.Render("(empty)"))
		w.WriteString(" ")
	}
	if commit.Description == "" {
		if commit.Empty {
			w.WriteString(s.Palette.EmptyPlaceholder.Render("(no description set)"))
		} else {
			w.WriteString(s.Palette.Placeholder.Render("(no description set)"))
		}
	} else {
		w.WriteString(s.Palette.Normal.Render(commit.Description))
	}
	return w.String()
}
