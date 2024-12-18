package test

import (
	"github.com/stretchr/testify/assert"
	"jjui/internal/jj"
	"os"
	"strings"
	"testing"
)

func Test_Parse_MergeTrees(t *testing.T) {
	testFiles := []string{
		//"testdata/merges.log",
		"testdata/merges-with-elided-revisions.rendered",
	}

	for _, fileName := range testFiles {
		fileName := fileName
		t.Run(fileName, func(t *testing.T) {
			file, err := os.Open(fileName)
			if err != nil {
				t.Fatalf("could not open file: %v", err)
			}

			p := jj.NewParser(file)
			lines := p.Parse()
			assert.NotEmpty(t, lines)
			assert.Len(t, lines, 10)
		})
	}
}

func Test_Parse_Tree(t *testing.T) {
	testFiles := []string{
		//"testdata/many-levels.log",
		//"testdata/elided-revisions.log",
		//"testdata/conflicted.log",
		//"testdata/merges.log",
		"testdata/merges-with-elided-revisions.log",
	}

	for _, fileName := range testFiles {
		fileName := fileName
		t.Run(fileName, func(t *testing.T) {
			file, err := os.Open(fileName)
			if err != nil {
				t.Fatalf("could not open file: %v", err)
			}

			p := jj.NewParser(file)
			var buffer strings.Builder
			rows := p.Parse()
			for _, row := range rows {
				jj.RenderRow(&buffer, row, TestRenderer{})
			}
			renderedFileName := strings.Replace(fileName, ".log", ".expected", 1)
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

func (t TestRenderer) RenderGlyph(connection jj.ConnectionType, commit *jj.Commit) string {
	return string(connection)
}

func (t TestRenderer) RenderTermination(connection jj.ConnectionType) string {
	return string(connection)
}

func (t TestRenderer) RenderChangeId(commit *jj.Commit) string {
	return commit.ChangeId
}

func (t TestRenderer) RenderAuthor(commit *jj.Commit) string {
	if commit.ChangeId == jj.RootChangeId {
		return "root()"
	}
	return ""
}

func (t TestRenderer) RenderDate(commit *jj.Commit) string {
	return ""
}

func (t TestRenderer) RenderBookmarks(commit *jj.Commit) string {
	return ""
}

func (t TestRenderer) RenderDescription(commit *jj.Commit) string {
	if commit.ChangeId == jj.RootChangeId {
		return ""
	}
	var w strings.Builder
	if commit.Empty {
		w.WriteString("(empty) ")
	}
	if commit.Description == "" {
		w.WriteString("(no description set)")
	} else {
		w.WriteString(commit.Description)
	}
	return w.String()
}
