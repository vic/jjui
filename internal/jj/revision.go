package jj

import (
	"bytes"
	"fmt"
	"os/exec"
)

const (
	TEMPLATE     = `separate(";", change_id.shortest(1), change_id.shortest(8), coalesce(bookmarks.join(","), "."), current_working_copy, immutable, conflict, empty, hidden, author.email(), author.timestamp().ago(), description.first_line())`
	RootChangeId = "zzzzzzzz"
)

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	IsWorkingCopy bool
	Author        string
	Timestamp     string
	Bookmarks     []string
	Description   string
	Immutable     bool
	Conflict      bool
	Empty         bool
	Hidden        bool
}

func (c Commit) IsRoot() bool {
	return c.ChangeId == RootChangeId
}

func (c Commit) GetChangeId() string {
	if c.Hidden {
		return "~" + c.ChangeId
	}
	return c.ChangeId
}

func (jj JJ) GetCommits(revset string) ([]GraphRow, error) {
	var args []string
	args = append(args, "log", "--color", "never", "--config", "ui.graph.style=curved", "--template", TEMPLATE)
	if revset != "" {
		args = append(args, "-r", revset)
	}
	cmd := exec.Command("jj", args...)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s", output)
	}
	p := NewParser(bytes.NewReader(output))
	graphLines := p.Parse()
	return graphLines, nil
}
