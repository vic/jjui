package test

import (
	"github.com/stretchr/testify/assert"
	"os"
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
	treeRenderer := jj.NewTreeRenderer(&dag)
	buffer := treeRenderer.RenderTree(TestRenderer{})
	content, err := os.ReadFile("testdata/many-levels.rendered")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	assert.Equal(t, string(content), buffer)
}

type TestRenderer struct{}

func (t TestRenderer) RenderCommit(commit *jj.Commit, context *jj.RenderContext) {
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
	context.Glyph = glyph
	context.Content = commit.ChangeIdShort
}

func (t TestRenderer) RenderElidedRevisions() string {
	return "~"
}
