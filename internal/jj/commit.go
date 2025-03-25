package jj

import "strings"

const (
	RootChangeId = "zzzzzzzz"
)

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	IsWorkingCopy bool
	Hidden        bool
	CommitIdShort string
	CommitId      string
}

func (c Commit) IsRoot() bool {
	return c.ChangeId == RootChangeId
}

func (c Commit) GetChangeId() string {
	if c.Hidden || strings.HasSuffix(c.ChangeId, "??") {
		return c.CommitId
	}
	return c.ChangeId
}
