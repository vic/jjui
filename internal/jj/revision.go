package jj

const (
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
	CommitIdShort string
	CommitId      string
}

func (c Commit) IsRoot() bool {
	return c.ChangeId == RootChangeId
}

func (c Commit) GetChangeId() string {
	if c.Hidden {
		return c.CommitId
	}
	return c.ChangeId
}
