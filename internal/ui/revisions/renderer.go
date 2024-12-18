package revisions

import (
	"fmt"
	"jjui/internal/jj"
	"strings"

	"jjui/internal/ui/common"

	tea "github.com/charmbracelet/bubbletea"
)

type SegmentedRenderer struct {
	Palette             common.Palette
	HighlightedRevision string
	Overlay             tea.Model
	op                  common.Operation
}

func (s SegmentedRenderer) RenderBefore(commit *jj.Commit) string {
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	if s.op == common.RebaseRevisionOperation || s.op == common.RebaseBranchOperation {
		if highlighted {
			return common.DropStyle.Render("<< here >>")
		}
	}
	return ""
}

func (s SegmentedRenderer) RenderAfter(commit *jj.Commit) string {
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	if highlighted && s.Overlay != nil && s.op != common.EditDescriptionOperation && s.op != common.SetBookmarkOperation {
		return s.Overlay.View()
	}
	return ""
}

func (s SegmentedRenderer) RenderGlyph(connection jj.ConnectionType, commit *jj.Commit) string {
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	style := s.Palette.Normal
	if highlighted {
		style = s.Palette.Selected
	}
	return style.Render(string(connection))
}

func (s SegmentedRenderer) RenderTermination(connection jj.ConnectionType) string {
	return s.Palette.CommitIdRestStyle.Render(string(connection))
}

func (s SegmentedRenderer) RenderChangeId(commit *jj.Commit) string {
	return fmt.Sprintf("%s%s", s.Palette.CommitShortStyle.Render(commit.ChangeIdShort), s.Palette.CommitIdRestStyle.Render(commit.ChangeId[len(commit.ChangeIdShort):]))
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
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	var w strings.Builder
	if s.op == common.SetBookmarkOperation && highlighted {
		w.WriteString(s.Overlay.View())
	}
	w.WriteString(s.Palette.BookmarksStyle.Render(strings.Join(commit.Bookmarks, " ")))
	return w.String()
}

func (s SegmentedRenderer) RenderDescription(commit *jj.Commit) string {
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	var w strings.Builder
	if s.op == common.EditDescriptionOperation && highlighted {
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
