package revisions

import (
	"fmt"
	"strings"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"

	tea "github.com/charmbracelet/bubbletea"
)

type SegmentedRenderer struct {
	Palette       common.Palette
	IsHighlighted bool
	Overlay       tea.Model
	op            common.Operation
	After         string
}

func (s SegmentedRenderer) RenderBefore(commit *jj.Commit) string {
	if s.IsHighlighted {
		if op, ok := s.op.(common.RebaseOperation); ok {
			if op.Target == common.RebaseTargetDestination {
				return fmt.Sprintf("%s %s onto %s", common.DropStyle.Render("<< onto >>"), s.Palette.CommitShortStyle.Render(op.From), s.Palette.CommitShortStyle.Render(commit.ChangeIdShort))
			}
			if op.Target == common.RebaseTargetAfter {
				return fmt.Sprintf("%s %s after %s", common.DropStyle.Render("<< after >>"), s.Palette.CommitShortStyle.Render(op.From), s.Palette.CommitShortStyle.Render(commit.ChangeIdShort))
			}
		}
	}
	return ""
}

func (s SegmentedRenderer) RenderAfter(commit *jj.Commit) string {
	if s.IsHighlighted && s.Overlay != nil {
		// TODO: neither of the following conditions should not match for overlay to be rendered
		if _, ok := s.op.(common.EditDescriptionOperation); !ok {
			return s.Overlay.View()
		}
		if _, ok := s.op.(common.SetBookmarkOperation); !ok {
			return s.Overlay.View()
		}
	}
	if s.IsHighlighted {
		if op, ok := s.op.(common.RebaseOperation); ok {
			if op.Target == common.RebaseTargetBefore {
				return fmt.Sprintf("%s %s before %s", common.DropStyle.Render("<< before >>"), s.Palette.CommitShortStyle.Render(op.From), s.Palette.CommitShortStyle.Render(commit.ChangeIdShort))
			}
		}
	}
	if s.After != "" {
		return s.After
	}
	return ""
}

func (s SegmentedRenderer) RenderGlyph(connection jj.ConnectionType, commit *jj.Commit) string {
	style := s.Palette.Normal
	squashDropMarker := ""
	if s.IsHighlighted {
		style = s.Palette.Selected
		if _, ok := s.op.(common.SquashOperation); ok {
			squashDropMarker = common.DropStyle.Render(" << into >> ")
		}
	}
	return style.Render(string(connection) + squashDropMarker)
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
	if _, ok := s.op.(common.SetBookmarkOperation); ok && s.IsHighlighted {
		w.WriteString(s.Overlay.View())
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
	if _, ok := s.op.(common.EditDescriptionOperation); ok && s.IsHighlighted {
		w.WriteString(s.Overlay.View())
		return w.String()
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
