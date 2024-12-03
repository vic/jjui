package revisions

import (
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

func (s *SegmentedRenderer) RenderCommit(commit *jj.Commit) string {
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	var w strings.Builder

	if (s.op == common.RebaseBranch || s.op == common.RebaseRevision) && highlighted {
		w.WriteString(common.DropStyle.Render("<< here >>"))
		w.WriteString("\n")
	}

	w.WriteString(s.Palette.CommitShortStyle.Render(commit.ChangeIdShort))
	w.WriteString(s.Palette.CommitIdRestStyle.Render(commit.ChangeId[len(commit.ChangeIdShort):]))
	w.Write([]byte{' '})

	w.WriteString(s.Palette.AuthorStyle.Render(commit.Author))
	w.Write([]byte{' '})

	w.WriteString(s.Palette.TimestampStyle.Render(commit.Timestamp))
	w.Write([]byte{' '})

	w.WriteString(s.Palette.BookmarksStyle.Render(strings.Join(commit.Bookmarks, " ")))
	w.Write([]byte{' '})

	if commit.Conflict {
		w.WriteString(s.Palette.ConflictStyle.Render("conflict"))
	}
	w.Write([]byte{'\n'})

	if commit.Empty {
		w.WriteString(s.Palette.Empty.Render("(empty)"))
		w.Write([]byte{' '})
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

	w.Write([]byte{'\n'})
	if s.Overlay != nil && highlighted {
		w.WriteString(s.Overlay.View())
	}
	return w.String()
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