package jj

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	TEMPLATE     = `separate(";", change_id.shortest(1), change_id.shortest(8), separate(",", parents.map(|x| x.change_id().shortest(1))), separate(",", coalesce(bookmarks, ".")), current_working_copy, immutable, conflict, empty, author.email(), author.timestamp().ago(), description.first_line())`
	RootChangeId = "zzzzzzzz"
)

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	Parents       []string
	IsWorkingCopy bool
	Author        string
	Timestamp     string
	Bookmarks     []string
	Description   string
	Immutable     bool
	Conflict      bool
	Empty         bool
}

func (c Commit) IsRoot() bool {
	return c.ChangeId == RootChangeId
}

func GetCommits(location string, revset string) (*Dag, error) {
	var args []string
	args = append(args, "log", "--reversed", "--template", TEMPLATE)
	if revset != "" {
		args = append(args, "-r", revset)
	}
	cmd := exec.Command("jj", args...)
	cmd.Dir = location
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s\n", output)
	}
	d := Parse(bytes.NewReader(output))
	return &d, nil
}

func Parse(reader io.Reader) Dag {
	d := NewDag()
	all, err := io.ReadAll(reader)
	if err != nil {
		return d
	}
	lines := strings.Split(string(all), "\n")
	stack := make([]*Node, 0)
	stack = append(stack, nil)
	levels := make([]int, 0)
	levels = append(levels, -1)
	seen := make(map[string]bool)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" || line == "~" {
			continue
		}
		index := strings.IndexAny(line, "○◆@×")
		if index == -1 {
			continue
		}
		_, after, _ := strings.Cut(line[index:], " ")
		parts := strings.Split(after, ";")
		commit := Commit{
			//TODO: change id can contain graph characters if there is a merge commit
			//this is a dirty hack to prevent panicking while trying to parse the graph.
			//parsed graph will be wrong but at least it won't panic until fixed.
			ChangeIdShort: strings.Trim(parts[0], " │"),
		}
		seen[commit.ChangeIdShort] = true
		if len(parts) > 1 {
			commit.ChangeId = parts[1]
		}
		edgeType := DirectEdge
		if len(parts) > 2 {
			commit.Parents = strings.Split(parts[2], ",")
			for _, parent := range commit.Parents {
				if _, ok := seen[parent]; !ok {
					edgeType = IndirectEdge
				}
			}
		}
		if len(parts) > 3 && parts[3] != "." {
			commit.Bookmarks = strings.Split(parts[3], ",")
		}
		if len(parts) > 4 {
			commit.IsWorkingCopy = parts[4] == "true"
		}
		if len(parts) > 5 {
			commit.Immutable = parts[5] == "true"
		}
		if len(parts) > 6 {
			commit.Conflict = parts[6] == "true"
		}
		if len(parts) > 7 {
			commit.Empty = parts[7] == "true"
		}
		if len(parts) > 8 {
			commit.Author = parts[8]
		}
		if len(parts) > 9 {
			commit.Timestamp = parts[9]
		}
		if len(parts) > 10 {
			commit.Description = parts[10]
		}
		node := d.AddNode(&commit)
		if index < levels[len(levels)-1] {
			levels = levels[:len(levels)-1]
			stack = stack[:len(stack)-1]
		}
		if stack[len(stack)-1] != nil {
			stack[len(stack)-1].AddEdge(node, edgeType)
		}
		if index == levels[len(levels)-1] {
			stack[len(stack)-1] = node
		}
		if index > levels[len(levels)-1] {
			levels = append(levels, index)
			stack = append(stack, node)
		}
		if commit.ChangeId == RootChangeId {
			commit.Conflict = false
			commit.Parents = nil
			commit.Immutable = false
			commit.Author = ""
			commit.Bookmarks = nil
			commit.Description = ""
		}
	}
	return d
}
