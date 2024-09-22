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

func TestRender_BranchedOrdered(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "zzwt", Index: 0, Parents: []string{"ulvv"}},
		{ChangeId: "ulvv", Index: 1, Parents: []string{"wzpt"}},
		{ChangeId: "wzpt", Index: 2, Parents: []string{"xumz"}},
		{ChangeId: "mxum", Index: 3, Parents: []string{"tklw"}},
		{ChangeId: "ywyr", Index: 4, Parents: []string{"tklw"}},
		{ChangeId: "mppl", Index: 5, Parents: []string{"tklw"}},
		{ChangeId: "tklw", Index: 6, Parents: []string{"qnor"}},
		{ChangeId: "xumz", Index: 7, Parents: []string{"mxsp"}},
		{ChangeId: "tnww", Index: 8, Parents: []string{"uolv"}},
		{ChangeId: "rnym", Index: 9, Parents: []string{"uolv"}},
		{ChangeId: "uolv", Index: 10, Parents: []string{"mvwv"}},
		{ChangeId: "ukkz", Index: 11, Parents: []string{"mpkz"}},
		{ChangeId: "mpkz", Index: 12, Parents: []string{"ssrp"}},
		{ChangeId: "puqn", Index: 13, Parents: []string{"ssrp"}},
		{ChangeId: "ssrp", Index: 14, Parents: []string{"vqrx"}},
		{ChangeId: "qnor", Index: 15, Parents: []string{"pnpu"}},
	}
	d := dag.Build(commits)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)

	expected := `
  ○ zzwt
  │ (no description)
  ○ ulvv
  │ (no description)
  ○ wzpt
  │ (no description)
  ○ xumz
  │ (no description)
  ~ (elided revisions)
  │ ○ tnww
  ├─╯ (no description)
  │ ○ rnym
  ├─╯ (no description)
  ○ uolv
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
