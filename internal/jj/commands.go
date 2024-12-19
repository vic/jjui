package jj

type Commands interface {
	GetConfig(key string) ([]byte, error)
	RebaseCommand(from string, to string) ([]byte, error)
	RebaseBranchCommand(from string, to string) ([]byte, error)
	SetDescription(rev string, description string) ([]byte, error)
	ListBookmark(revision string) ([]string, error)
	SetBookmark(revision string, name string) ([]byte, error)
	MoveBookmark(revision string, bookmark string) ([]byte, error)
	DeleteBookmark(bookmark string) ([]byte, error)
	GitFetch() ([]byte, error)
	GitPush() ([]byte, error)
	Diff(revision string) ([]byte, error)
	Edit(revision string) ([]byte, error)
	DiffEdit(revision string) ([]byte, error)
	Abandon(revision string) ([]byte, error)
	New(from string) ([]byte, error)
	GetCommits(revset string) ([]GraphRow, error)
}
