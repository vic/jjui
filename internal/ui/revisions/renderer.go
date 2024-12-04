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

func (s *SegmentedRenderer) RenderCommit(commit *jj.Commit, context *jj.RenderContext) {
	highlighted := commit.ChangeIdShort == s.HighlightedRevision
	if (s.op == common.RebaseBranch || s.op == common.RebaseRevision) && highlighted {
		context.RenderLine(common.DropStyle.Render("<< here >>"))
	}

	style := s.Palette.Normal
	if highlighted {
		style = s.Palette.Selected
	}
	if commit.Immutable {
		context.SetGlyph(style.Render("◆"))
	} else if commit.IsWorkingCopy {
		context.SetGlyph(style.Render("@"))
	} else {
		context.SetGlyph(style.Render("○"))
	}

	var w strings.Builder
	w.WriteString(s.Palette.CommitShortStyle.Render(commit.ChangeIdShort))
	w.WriteString(s.Palette.CommitIdRestStyle.Render(commit.ChangeId[len(commit.ChangeIdShort):]))
	w.WriteString(" ")

	w.WriteString(s.Palette.AuthorStyle.Render(commit.Author))
	w.WriteString(" ")

	w.WriteString(s.Palette.TimestampStyle.Render(commit.Timestamp))
	w.WriteString(" ")

	w.WriteString(s.Palette.BookmarksStyle.Render(strings.Join(commit.Bookmarks, " ")))

	if commit.Conflict {
		w.WriteString(" ")
		w.WriteString(s.Palette.ConflictStyle.Render("conflict"))
	}
	context.RenderLine(w.String())

	w.Reset()
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

	context.RenderLine(w.String())
	if s.Overlay != nil && highlighted {
		context.RenderLine(s.Overlay.View())
	}
}

func (s *SegmentedRenderer) RenderElidedRevisions() string {
	return s.Palette.CommitIdRestStyle.Render("~  (elided revisions)")
}
