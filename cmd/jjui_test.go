package main

import (
	"jjui/internal/dag"
	"jjui/internal/jj"
	"jjui/internal/ui/revisions"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	model := revisions.New()
	commit := jj.Commit{
		ChangeIdShort: "top",
		ChangeId:      "topchange",
		Parents:       nil,
	}
	root := dag.Build([]jj.Commit{commit})
	model.rows = dag.BuildGraphRows(root)

	view := model.View()
	expected := `○ topchange 
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
