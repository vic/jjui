package jj

import (
	"bufio"
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type GraphWriter struct {
	buffer             bytes.Buffer
	lineCount          int
	connectionPos      int
	connections        []ConnectionType
	connectionsWritten bool
	renderer           RowRenderer
	row                GraphLine
}

func (w *GraphWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if !w.connectionsWritten {
		w.renderConnections()
	}

	if bytes.Index(p, []byte("\n")) == -1 {
		return w.buffer.Write(p)
	}
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		line := scanner.Text()
		if !w.connectionsWritten {
			w.renderConnections()
		}
		w.buffer.Write([]byte(line))
		w.buffer.Write([]byte("\n"))
		w.lineCount++
		w.connectionsWritten = false
		w.connections = extendConnections(w.connections)
	}

	return n, nil
}

func (w *GraphWriter) LineCount() int {
	return w.lineCount
}

func (w *GraphWriter) String(start, end int) string {
	lines := strings.Split(w.buffer.String(), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	if start < 0 {
		start = 0
	}
	h := end - start
	for h > len(lines) {
		lines = append(lines, "")
	}
	return strings.Join(lines[start:end], "\n")
}

func (w *GraphWriter) Reset() {
	w.buffer.Reset()
	w.lineCount = 0
}

func (w *GraphWriter) RenderRow(row GraphLine, renderer RowRenderer) {
	w.connectionPos = 0
	w.connectionsWritten = false
	w.row = row
	w.renderer = renderer
	w.connections = extendConnections(w.connections)
	// will render by extending the previous connections
	written, _ := w.Write([]byte(renderer.RenderBefore(row.Commit)))
	if written > 0 {
		w.Write([]byte("\n"))
	}
	w.connectionsWritten = false
	w.connections = row.Connections[0]
	fmt.Fprintf(w, renderer.RenderChangeId(row.Commit))
	if author := renderer.RenderAuthor(row.Commit); author != "" {
		fmt.Fprintf(w, " ")
		fmt.Fprintf(w, author)
	}
	if date := renderer.RenderDate(row.Commit); date != "" {
		fmt.Fprintf(w, " ")
		fmt.Fprintf(w, date)
	}
	if bookmarks := renderer.RenderBookmarks(row.Commit); bookmarks != "" {
		fmt.Fprintf(w, " ")
		fmt.Fprintf(w, bookmarks)
	}
	fmt.Fprintln(w)

	if row.Commit.IsRoot() {
		return
	}

	lastLineConnection := extendConnections(row.Connections[0])
	if len(row.Connections) > 1 && !slices.Contains(row.Connections[1], TERMINATION) {
		w.connectionPos = 1
		lastLineConnection = row.Connections[1]
	}

	if description := renderer.RenderDescription(row.Commit); description != "" {
		lines := strings.Split(description, "\n")
		n := len(lines)
		for i, line := range lines {
			if i == n-1 {
				w.connections = lastLineConnection
			} else {
				w.connections = extendConnections(row.Connections[0])
			}
			fmt.Fprintf(w, line)
			fmt.Fprintln(w)
		}
	}

	w.connections = extendConnections(lastLineConnection)
	written, _ = w.Write([]byte(renderer.RenderAfter(row.Commit)))
	if written > 0 {
		w.Write([]byte("\n"))
	}

	w.connectionPos++
	for w.connectionPos < len(row.Connections) {
		w.connections = row.Connections[w.connectionPos]
		w.renderConnections()
		if slices.Contains(w.connections, TERMINATION) {
			w.buffer.Write([]byte(w.renderer.RenderTermination("(elided revisions)")))
		}
		w.buffer.Write([]byte("\n"))
		w.lineCount++
		w.connectionPos++
	}
}

func (w *GraphWriter) renderConnections() {
	if w.connections == nil {
		w.connectionsWritten = true
		return
	}
	maxPadding := 0
	for _, c := range w.row.Connections {
		if len(c) > maxPadding {
			maxPadding = len(c)
		}
	}

	for _, c := range w.connections {
		if c == GLYPH || c == GLYPH_IMMUTABLE || c == GLYPH_WORKING_COPY || c == GLYPH_CONFLICT {
			w.buffer.WriteString(w.renderer.RenderGlyph(c, w.row.Commit))
		} else if c == TERMINATION {
			w.buffer.WriteString(w.renderer.RenderTermination(c))
		} else {
			w.buffer.WriteString(string(c))
		}
	}
	if len(w.connections) < maxPadding {
		w.buffer.WriteString(strings.Repeat(SPACE, maxPadding-len(w.connections)))
	}
	w.buffer.WriteString(" ")
	w.connectionsWritten = true
}

func extendConnections(connections []ConnectionType) []ConnectionType {
	if connections == nil {
		return nil
	}
	extended := make([]ConnectionType, 0)
	for _, cur := range connections {
		if cur != MERGE_LEFT && cur != MERGE_BOTH && cur != MERGE_RIGHT {
			extended = append(extended, VERTICAL)
		}
	}
	return extended
}
