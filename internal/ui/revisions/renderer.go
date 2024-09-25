package revisions

import (
    "strings"

    "jjui/internal/dag"
    "jjui/internal/ui/common"
)

type (
	ChangeIdShort   struct{}
	ChangeIdRest    struct{}
	Author          struct{}
	Branches        struct{}
	ConflictMarker  struct{}
	Description     struct{}
	NodeGlyph       struct{}
	Glyph           struct{}
	ElidedRevisions struct{}
)

func SegmentedRenderer(w *strings.Builder, row *dag.GraphRow, palette common.Palette, highlighted bool, segments ...interface{}) {
	for _, segment := range segments {
		switch segment := segment.(type) {
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
		case ChangeIdShort:
			w.WriteString(palette.CommitShortStyle.Render(row.Commit.ChangeIdShort))
		case ChangeIdRest:
			w.WriteString(palette.CommitIdRestStyle.Render(row.Commit.ChangeId[len(row.Commit.ChangeIdShort):]))
		case Author:
			w.WriteString(palette.AuthorStyle.Render(row.Commit.Author))
		case Branches:
			w.WriteString(palette.BranchesStyle.Render(row.Commit.Branches))
		case ConflictMarker:
			if row.Commit.Conflict {
				w.WriteString(palette.ConflictStyle.Render("conflict"))
			}
		case Description:
			if row.Commit.Description == "" {
				if row.Commit.Empty {
					w.WriteString(palette.Empty.Render("(empty) (no description)"))
				} else {
					w.WriteString(palette.NonEmpty.Render("(no description)"))
				}
			} else {
				w.WriteString(palette.Normal.Render(row.Commit.Description))
			}

		case string:
			w.WriteString(segment)
		}
	}
}
