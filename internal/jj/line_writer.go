package jj

import (
	"bytes"
	"strings"
)

// LineTrackingWriter is a custom writer that tracks lines and stores content
type LineTrackingWriter struct {
	buffer    bytes.Buffer
	lineCount int
}

func (w *LineTrackingWriter) Write(p []byte) (n int, err error) {
	// Write to the buffer
	n, err = w.buffer.Write(p)
	if err != nil {
		return n, err
	}

	if bytes.Index(p, []byte("\n")) == -1 {
		return n, nil
	}

	// Split the written content into lines
	newLines := strings.Split(string(p), "\n")
	for _, line := range newLines {
		if line != "" {
			w.lineCount++
		}
	}

	return n, nil
}

func (w *LineTrackingWriter) LineCount() int {
	return w.lineCount
}

func (w *LineTrackingWriter) String(start, end int) string {
	lines := strings.Split(w.buffer.String(), "\n")
	if start < 0 {
		start = 0
	}
	h := end - start
	for h > len(lines) {
		lines = append(lines, "")
	}
	return strings.Join(lines[start:end], "\n")
}

func (w *LineTrackingWriter) Reset() {
	w.buffer.Reset()
	w.lineCount = 0
}
