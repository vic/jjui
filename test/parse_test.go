package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"

	"jjui/internal/dag"
	"jjui/internal/jj"
)

func parse(reader io.Reader) *dag.Dag {
	all, err := io.ReadAll(reader)
	if err != nil {
		return nil
	}
	d := dag.NewDag()
	lines := strings.Split(string(all), "\n")
	stack := make([]*dag.Node, 0)
	stack = append(stack, nil)
	levels := make([]int, 0)
	levels = append(levels, -1)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" || line == "~" {
			continue
		}
		index := strings.IndexAny(line, "○◆@")
		if index == -1 {
			continue
		}
		_, after, _ := strings.Cut(line[index:], " ")
		parts := strings.Split(after, ";")
		commit := jj.Commit{
			ChangeIdShort: strings.TrimSpace(parts[0]),
		}
		if len(parts) > 1 {
			commit.ChangeId = parts[1]
		}
		if len(parts) > 2 {
			commit.IsWorkingCopy = parts[2] == "true"
		}
		if len(parts) > 3 {
			commit.Immutable = parts[3] == "true"
		}
		if len(parts) > 4 {
			commit.Conflict = parts[4] == "true"
		}
		if len(parts) > 5 {
			commit.Empty = parts[5] == "true"
		}
		if len(parts) > 6 {
			commit.Author = parts[6]
		}
		if len(parts) > 7 {
			commit.Timestamp = parts[7]
		}
		if len(parts) > 8 {
			commit.Description = parts[8]
		}
		node := d.AddNode(&commit)
		if index < levels[len(levels)-1] {
			levels = levels[:len(levels)-1]
			stack = stack[:len(stack)-1]
		}
		if stack[len(stack)-1] != nil {
			stack[len(stack)-1].AddEdge(node, dag.DirectEdge)
		}
		if index == levels[len(levels)-1] {
			stack[len(stack)-1] = node
		}
		if index > levels[len(levels)-1] {
			levels = append(levels, index)
			stack = append(stack, node)
		}
	}
	rows := dag.BuildGraphRows(d.GetRoot())
	fmt.Printf("%v\n", rows)
	return d
}

func Test_parseLogOutput_ManyLevels(t *testing.T) {
	fileName := "testdata/many-levels.log"
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	d := parse(file)
	assert.Equal(t, 8, len(d.Nodes))
}

func Test_parseLogOutput_TwoLevels(t *testing.T) {
	fileName := "testdata/two-level.log"
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	d := parse(file)
	assert.Equal(t, 10, len(d.Nodes))
}
