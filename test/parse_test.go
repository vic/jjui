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
	treeRenderer := jj.NewTreeRenderer(&dag, TestRenderer{})
	buffer := treeRenderer.RenderTree()
	content, err := os.ReadFile("testdata/many-levels.rendered")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	assert.Equal(t, string(content), buffer)
}

type TestRenderer struct { }

func (t TestRenderer) RenderCommit(commit *jj.Commit) string {
	return commit.ChangeIdShort
}

func (t TestRenderer) RenderElidedRevisions() string {
	//TODO implement me
	return "~  (elided revisions)"
}

func (t TestRenderer) RenderGlyph(commit *jj.Commit) string {
	if commit.Immutable {
		return "◆  "
	} else if commit.IsWorkingCopy {
		return "@  "
	} else {
		return "○  "
	}
}
