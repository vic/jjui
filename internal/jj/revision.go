package jj

import (
	"bytes"
	"fmt"
	"os/exec"
)

const (
	TEMPLATE     = `separate(";", change_id.shortest(1), change_id.shortest(8), separate(",", coalesce(bookmarks, ".")), current_working_copy, immutable, conflict, empty, author.email(), author.timestamp().ago(), description.first_line())`
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

func GetCommits(location string, revset string) ([]GraphLine, error) {
	var args []string
	args = append(args, "log", "--template", TEMPLATE)
	if revset != "" {
		args = append(args, "-r", revset)
	}
	cmd := exec.Command("jj", args...)
	cmd.Dir = location
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s\n", output)
	}
	p := NewParser(bytes.NewReader(output))
	graphLines := p.Parse()
	return graphLines, nil
}
