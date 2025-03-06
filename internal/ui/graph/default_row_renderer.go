package graph

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
	"strings"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/operations"
)

type DefaultRowRenderer struct {
	Palette       common.Palette
	IsHighlighted bool
	Op            operations.Operation
}

func (s DefaultRowRenderer) RenderConnection(connectionType jj.ConnectionType) string {
	return s.Palette.Normal.Render(string(connectionType))
}

func (s DefaultRowRenderer) RenderBefore(*jj.Commit) string {
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionBefore {
		return s.Op.Render()
	}
	return ""
}

func (s DefaultRowRenderer) RenderAfter(*jj.Commit) string {
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionAfter {
		return s.Op.Render()
	}
	return ""
}

func (s DefaultRowRenderer) RenderGlyph(connection jj.ConnectionType, commit *jj.Commit) string {
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

func (s DefaultRowRenderer) RenderTermination(connection jj.ConnectionType) string {
	return s.Palette.Elided.Render(string(connection))
}

func (s DefaultRowRenderer) RenderChangeId(commit *jj.Commit) string {
	changeId := s.Palette.ChangeId.Render("", commit.ChangeIdShort) + s.Palette.Rest.Render(commit.ChangeId[len(commit.ChangeIdShort):])
	hidden := ""
	if commit.Hidden {
		hidden = s.Palette.Normal.Render(" hidden")
		return lipgloss.JoinHorizontal(0, changeId, hidden)
	}
	return changeId
}

func (s DefaultRowRenderer) RenderCommitId(commit *jj.Commit) string {
	if commit.IsRoot() {
		return ""
	}
	return s.Palette.CommitId.Render("", commit.CommitIdShort) + s.Palette.Rest.Render(commit.CommitId[len(commit.ChangeIdShort):])
}

func (s DefaultRowRenderer) RenderAuthor(commit *jj.Commit) string {
	if commit.IsRoot() {
		return s.Palette.EmptyPlaceholder.Render("root()")
	}
	return s.Palette.Author.Render("", commit.Author)
}

func (s DefaultRowRenderer) RenderDate(commit *jj.Commit) string {
	if commit.IsRoot() {
		return ""
	}
	return s.Palette.Timestamp.Render("", commit.Timestamp)
}

func (s DefaultRowRenderer) RenderBookmarks(commit *jj.Commit) string {
	var w strings.Builder
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionBookmark {
		w.WriteString(" ")
		w.WriteString(s.Op.Render())
	}
	if len(commit.Bookmarks) > 0 {
		var bookmarks []string
		bookmarks = append(bookmarks, "")
		bookmarks = append(bookmarks, commit.Bookmarks...)
		w.WriteString(s.Palette.Bookmarks.Render(bookmarks...))
	}
	return w.String()
}

func (s DefaultRowRenderer) RenderMarkers(commit *jj.Commit) string {
	if commit.Conflict {
		return s.Palette.Conflict.Render("conflict")
	}
	return ""
}

func (s DefaultRowRenderer) RenderDescription(commit *jj.Commit) string {
	if s.IsHighlighted && s.Op.RenderPosition() == operations.RenderPositionDescription {
		return s.Op.Render()
	}

	if commit.Empty && commit.Description == "" {
		return s.Palette.EmptyPlaceholder.Render(" (empty) (no description set)")
	}
	if commit.Empty {
		return lipgloss.JoinHorizontal(0, s.Palette.EmptyPlaceholder.Render(" (empty)", s.Palette.Normal.Render(commit.Description)))
	} else if commit.Description == "" {
		return s.Palette.Placeholder.Render(" (no description set)")
	}
	return s.Palette.Normal.Render("", commit.Description)
}
