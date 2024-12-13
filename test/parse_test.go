package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"

	"jjui/internal/jj"
)

func Test_Parse_Tree(t *testing.T) {
	fileName := "testdata/many-levels.log"
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer file.Close()

	dag := jj.Parse(file)
	var buffer strings.Builder
	rows := dag.GetTreeRows()
	for _, row := range rows {
		jj.RenderRow(&buffer, row, TestRenderer{})
	}
	content, err := os.ReadFile("testdata/many-levels.rendered")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	assert.Equal(t, string(content), buffer.String())
}

type TestRenderer struct{}

func (t TestRenderer) Render(row *jj.TreeRow) {
	commit := row.Commit
	glyph := ""
	if commit.Immutable {
		glyph = "◆"
	} else if commit.IsWorkingCopy {
		glyph = "@"
	} else if commit.Conflict {
		glyph = "×"
	} else {
		glyph = "○"
	}
	row.Glyph = glyph
	row.Content = commit.ChangeIdShort
	if row.EdgeType == jj.IndirectEdge {
		row.ElidedRevision = "~"
	}
}
