package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"jjui/internal/dag"
	"jjui/internal/jj"
	"jjui/internal/ui/revisions"

	"github.com/stretchr/testify/assert"
)

func TestRender_Single(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "topchange"},
	}
	parents := make(map[string]string)
	parents["topchange"] = ""
	d := dag.Build(commits, parents)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)
	model, _ = model.Update(tea.WindowSizeMsg{
		Width:  100,
		Height: 100,
	})

	expected := `○ topchange
                 │ (no description set)`
	verifyOutput(t, expected, model.View())
}

func TestRender_ElidedRevisions(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "a", Parents: []string{"c"}},
		{ChangeId: "b"},
	}
	parents := make(map[string]string)
	parents["a"] = "c"
	parents["c"] = "b"
	parents["b"] = ""
	d := dag.Build(commits, parents)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)
	model, _ = model.Update(tea.WindowSizeMsg{
		Width:  100,
		Height: 100,
	})

	expected := `
  ○ a  
  │ (no description set)
  ~ (elided revisions)
  ○ b  
  │ (no description set)
  `
	verifyOutput(t, expected, model.View())
}

func TestRender_Branched(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "a1", Parents: []string{"root"}},
		{ChangeId: "b", Parents: []string{"root"}},
		{ChangeId: "root"},
	}
	parents := make(map[string]string)
	parents["a1"] = "root"
	parents["b"] = "root"
	parents["root"] = ""
	d := dag.Build(commits, parents)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)
	model, _ = model.Update(tea.WindowSizeMsg{
		Width:  100,
		Height: 100,
	})

	expected := `
  ○ a1
  │ (no description set)
  │ ○ b
  ├─╯ (no description set)
  ○ root
  │ (no description set)
  `
	verifyOutput(t, expected, model.View())
}

func TestRender_BranchedOrdered(t *testing.T) {
	commits := []jj.Commit{
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
	parents := make(map[string]string)
	parents["mxsp"] = "uolv"
	parents["vqrx"] = "qnor"

	parents["wzpt"] = "xumz"
	parents["mxum"] = "tklw"
	parents["ywyr"] = "tklw"
	parents["mppl"] = "tklw"
	parents["tklw"] = "qnor"
	parents["xumz"] = "mxsp"
	parents["tnww"] = "uolv"
	parents["rnym"] = "uolv"
	parents["uolv"] = "mvwv"
	parents["ukkz"] = "mpkz"
	parents["mpkz"] = "ssrp"
	parents["puqn"] = "ssrp"
	parents["ssrp"] = "vqrx"
	parents["qnor"] = "pnpu"

	d := dag.Build(commits, parents)
	root := dag.BuildGraphRows(d)
	model := revisions.New(root)
	model, _ = model.Update(tea.WindowSizeMsg{
		Width:  100,
		Height: 100,
	})

	expected := `
  ○ wzpt
  │ (no description set)
  ○ xumz
  │ (no description set)
  ~ (elided revisions)
  │ ○ tnww
  ├─╯ (no description set)
  │ ○ rnym
  ├─╯ (no description set)
  ○ uolv
  │ (no description set)
  `
	verifyOutput(t, expected, model.View())
}

func verifyOutput(t *testing.T, expected, view string) {
	expected = deindent(expected)
	actual := deindent(view)
	assert.Equal(t, expected, actual)
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
	ret := strings.Join(output, "\n")
	return strings.TrimRight(ret, "\n")
}
