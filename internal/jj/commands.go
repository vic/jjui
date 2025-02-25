package jj

type Commands interface {
	GetConfig(key string) ([]byte, error)
	RebaseCommand(from string, to string, source string, target string) Command
	SetDescription(rev string, description string) Command
	ListBookmark(revision string) Command
	SetBookmark(revision string, name string) Command
	MoveBookmark(revision string, bookmark string) Command
	DeleteBookmark(bookmark string) Command
	GitFetch() Command
	GitPush() Command
	Diff(revision string, fineName string) Command
	Edit(revision string) Command
	DiffEdit(revision string) Command
	Abandon(revision string) Command
	New(from string) Command
	Split(revision string, files []string) Command
	GetCommits(revset string) ([]GraphRow, error)
	Squash(from string, destination string) Command
	Status(revision string) Command
	Restore(revision string, files []string) Command
	Undo() Command
	Show(revision string) Command
}
