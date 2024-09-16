package main

import (
	"github.com/stretchr/testify/assert"
	"jjui/internal/dag"
	"jjui/internal/jj"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	model := initialModel()
	commit := jj.Commit{
		ChangeIdShort: "top",
		ChangeId:      "topchange",
		Parents:       nil,
	}
	root := dag.Build([]jj.Commit{commit})
	model.rows = dag.BuildGraphRows(root)

	view := model.View()
	expected := `○ topchange  edges: 0 level: 0 
                 │ (no description)
                 use j,k keys to move up and down: cursor:0 dragged:-1`
	assert.Equal(t, deindent(expected), view)
}

func deindent(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimLeft(lines[i], " \t")
	}
	return strings.Join(lines, "\n") + "\n"
}
