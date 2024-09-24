package revisions

import (
	"strings"

	"jjui/internal/dag"
	"jjui/internal/ui/common"
)

func DefaultRenderer(w *strings.Builder, row *dag.GraphRow, palette common.Palette, highlighted bool) {
	indent := strings.Repeat("│ ", row.Level)
	glyph := "│"
	nodeGlyph := "○ "
	if row.Commit.IsWorkingCopy {
		nodeGlyph = "@ "
	}
	if row.Commit.Immutable {
		nodeGlyph = "◆ "
	}

	if !row.IsFirstChild {
		indent = strings.Repeat("│ ", row.Level-1)
		glyph = "├─╯"
		nodeGlyph = "│ " + nodeGlyph
	}
	w.WriteString(indent)
	if highlighted {
		w.WriteString(palette.Selected.Render(nodeGlyph))
	} else {
		w.WriteString(nodeGlyph)
	}
	w.WriteString(palette.CommitShortStyle.Render(row.Commit.ChangeIdShort))
	w.WriteString(palette.CommitIdRestStyle.Render(row.Commit.ChangeId[len(row.Commit.ChangeIdShort):]))
	w.WriteString(" ")
	w.WriteString(palette.AuthorStyle.Render(row.Commit.Author))
	w.WriteString(" ")
	w.WriteString(palette.BranchesStyle.Render(row.Commit.Branches))
	if row.Commit.Conflict {
		w.WriteString(" ")
		w.WriteString(palette.ConflictStyle.Render("conflict"))
	}
	w.WriteString("\n")
	// description line
	w.WriteString(indent)
	w.WriteString(glyph)
	w.WriteString(" ")
	if row.Commit.Description == "" {
		if row.Commit.Empty {
			w.WriteString(palette.Empty.Render("(empty) (no description)"))
		} else {
			w.WriteString(palette.NonEmpty.Render("(no description)"))
		}
	} else {
		w.WriteString(palette.Normal.Render(row.Commit.Description))
	}
	w.WriteString("\n")
	if row.Elided {
		w.WriteString(indent)
		w.WriteString(palette.CommitIdRestStyle.Render("~ (elided revisions)"))
		w.WriteString("\n")
	}
}
