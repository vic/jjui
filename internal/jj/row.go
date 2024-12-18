package jj

import (
	"fmt"
	"io"
	"slices"
)

type RowRenderer interface {
	RenderBefore(commit *Commit) string
	RenderAfter(commit *Commit) string
	RenderGlyph(connection ConnectionType, commit *Commit) string
	RenderTermination(connection ConnectionType) string
	RenderChangeId(commit *Commit) string
	RenderAuthor(commit *Commit) string
	RenderDate(commit *Commit) string
	RenderBookmarks(commit *Commit) string
	RenderDescription(commit *Commit) string
}

func RenderRow(w io.Writer, row GraphLine, renderer RowRenderer) {
	for i, connections := range ensureHeight(row.Connections, 2) {
		renderElidedRevisions := false
		for _, connection := range connections {
			if connection == GLYPH || connection == GLYPH_IMMUTABLE || connection == GLYPH_WORKING_COPY || connection == GLYPH_CONFLICT {
				fmt.Fprintf(w, renderer.RenderGlyph(connection, row.Commit))
			} else if connection == TERMINATION {
				renderElidedRevisions = true
				fmt.Fprintf(w, renderer.RenderTermination(connection))
			} else {
				fmt.Fprintf(w, string(connection))
			}
		}

		if i == 0 {
			fmt.Fprintf(w, " ")
			written, _ := fmt.Fprintf(w, renderer.RenderChangeId(row.Commit))
			if written > 0 {
				fmt.Fprintf(w, " ")
			}
			written, _ = fmt.Fprintf(w, renderer.RenderAuthor(row.Commit))
			if written > 0 {
				fmt.Fprintf(w, " ")
			}
			written, _ = fmt.Fprintf(w, renderer.RenderDate(row.Commit))
			if written > 0 {
				fmt.Fprintf(w, " ")
			}
			written, _ = fmt.Fprintf(w, renderer.RenderBookmarks(row.Commit))
		}
		if row.Commit.IsRoot() {
			fmt.Fprintln(w)
			break
		}
		if i == 1 {
			fmt.Fprintf(w, " ")
			fmt.Fprintf(w, renderer.RenderDescription(row.Commit))
		}
		if renderElidedRevisions {
			fmt.Fprintf(w, " ")
			fmt.Fprintf(w, renderer.RenderTermination("(elided revisions)"))
		}
		fmt.Fprintln(w)
	}
}

func extendConnection(connections []ConnectionType) []ConnectionType {
	extended := make([]ConnectionType, 0)
	for _, cur := range connections {
		if cur == GLYPH || cur == GLYPH_IMMUTABLE || cur == GLYPH_WORKING_COPY || cur == GLYPH_CONFLICT {
			extended = append(extended, VERTICAL)
		} else if cur == TERMINATION {
			extended = append(extended, VERTICAL)
		} else if cur == VERTICAL {
			extended = append(extended, VERTICAL)
		} else if cur == JOIN_RIGHT || cur == JOIN_BOTH || cur == JOIN_LEFT {
			extended = append(extended, VERTICAL)
		} else {
			extended = append(extended, VERTICAL)
		}
	}
	return extended
}

func ensureHeight(connections [][]ConnectionType, count int) [][]ConnectionType {
	availableSpace := 0
	extended := make([][]ConnectionType, 0)
	for _, connection := range connections {
		if slices.Contains(connection, TERMINATION) {
			for availableSpace < count {
				extended = append(extended, extendConnection(connection))
				availableSpace++
			}
		}
		availableSpace++
		extended = append(extended, connection)
	}
	for len(extended) < count {
		extended = append(extended, extendConnection(extended[len(extended)-1]))
	}
	return extended
}
