package main

import (
	"strings"
	"testing"

	"jjui/internal/dag"
	"jjui/internal/jj"
	"jjui/internal/ui/revisions"

	"github.com/stretchr/testify/assert"
)

func TestRender_Single(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "topchange"},
	}
	d := dag.Build(commits)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)

	expected := `○ topchange
               │ (no description)`
	verifyOutput(t, expected, model.View())
}

func TestRender_ElidedRevisions(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "a"},
		{ChangeId: "b"},
	}
	d := dag.Build(commits)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)

	expected := `
  ○ a  
  │ (no description)
  ~ (elided revisions)
  ○ b  
  │ (no description)
  `
	verifyOutput(t, expected, model.View())
}

func TestRender_Branched(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "a1", Parents: []string{"root"}},
		{ChangeId: "b", Parents: []string{"root"}},
		{ChangeId: "root"},
	}
	d := dag.Build(commits)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)

	expected := `
  ○ a1
  │ (no description)
  │ ○ b
  ├─╯ (no description)
  ○ root
  │ (no description)
  `
	verifyOutput(t, expected, model.View())
}

func verifyOutput(t *testing.T, expected, view string) {
	expected = deindent(expected)
	actual := deindent(view)
	prefix := actual[:len(expected)]
	assert.Equal(t, expected, prefix)
}

func deindent(s string) string {
	lines := strings.Split(s, "\n")
	var output []string
	for i := range lines {
		line := lines[i]
		if line == "" {
			continue
		}
		line = strings.TrimLeft(line, " \t")
		line = strings.TrimSpace(line)
		output = append(output, line)
	}
	return strings.Join(output, "\n")
}
