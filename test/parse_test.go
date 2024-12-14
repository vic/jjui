package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"

	"jjui/internal/jj"
)

func Test_Parse_Tree(t *testing.T) {
	testFiles := []string{
		"testdata/many-levels.log",
		"testdata/elided-revisions.log",
		"testdata/conflicted.log",
	}

	for _, fileName := range testFiles {
		fileName := fileName
		t.Run(fileName, func(t *testing.T) {
			file, err := os.Open(fileName)
			if err != nil {
				t.Fatalf("could not open file: %v", err)
			}

			dag := jj.Parse(file)
			var buffer strings.Builder
			rows := dag.GetTreeRows()
			for _, row := range rows {
				jj.RenderRow(&buffer, row, TestRenderer{})
			}
			renderedFileName := strings.Replace(fileName, ".log", ".rendered", 1)
			content, err := os.ReadFile(renderedFileName)
			if err != nil {
				t.Fatalf("could not read file: %v", err)
			}
			_ = file.Close()
			assert.Equal(t, string(content), buffer.String())
		})
	}
}

type TestRenderer struct{}

func (t TestRenderer) Render(row *jj.TreeRow) {
	commit := row.Commit
	glyph := ""
	if commit.Immutable || commit.IsRoot() {
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
