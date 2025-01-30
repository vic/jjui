package jj

import "os/exec"

type Commands interface {
	GetConfig(key string) ([]byte, error)
	RebaseCommand(from string, to string) *exec.Cmd
	RebaseBranchCommand(from string, to string) *exec.Cmd
	SetDescription(rev string, description string) *exec.Cmd
	ListBookmark(revision string) *exec.Cmd
	SetBookmark(revision string, name string) *exec.Cmd
	MoveBookmark(revision string, bookmark string) *exec.Cmd
	DeleteBookmark(bookmark string) *exec.Cmd
	GitFetch() *exec.Cmd
	GitPush() *exec.Cmd
	Diff(revision string, fineName string) *exec.Cmd
	Edit(revision string) *exec.Cmd
	DiffEdit(revision string) *exec.Cmd
	Abandon(revision string) *exec.Cmd
	New(from string) *exec.Cmd
	Split(revision string) *exec.Cmd
	GetCommits(revset string) ([]GraphRow, error)
	Squash(from string, destination string) *exec.Cmd
	Status(revision string) *exec.Cmd
	Restore(revision string, files []string) *exec.Cmd
}
