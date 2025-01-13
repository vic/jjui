package jj

import (
	"bytes"
	"fmt"
	"os/exec"
)

type JJ struct {
	Location string
}

func (jj JJ) GetConfig(key string) ([]byte, error) {
	cmd := exec.Command("jj", "config", "get", key)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	output = bytes.Trim(output, "\n")
	return output, err
}

func (jj JJ) RebaseCommand(from string, to string) *exec.Cmd {
	cmd := exec.Command("jj", "rebase", "-r", from, "-d", to)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) RebaseBranchCommand(from string, to string) *exec.Cmd {
	cmd := exec.Command("jj", "rebase", "-b", from, "-d", to)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) SetDescription(rev string, description string) *exec.Cmd {
	cmd := exec.Command("jj", "describe", "-r", rev, "-m", description)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) ListBookmark(revision string) *exec.Cmd {
	cmd := exec.Command("jj", "log", "-r", fmt.Sprintf("::%s- & bookmarks()", revision), "--template", "local_bookmarks.map(|x| x.name() ++ '\n')", "--no-graph")
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) SetBookmark(revision string, name string) *exec.Cmd {
	cmd := exec.Command("jj", "bookmark", "set", "-r", revision, name)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) MoveBookmark(revision string, bookmark string) *exec.Cmd {
	cmd := exec.Command("jj", "bookmark", "move", bookmark, "--to", revision)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) DeleteBookmark(bookmark string) *exec.Cmd {
	cmd := exec.Command("jj", "bookmark", "delete", bookmark)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) GitFetch() *exec.Cmd {
	cmd := exec.Command("jj", "git", "fetch")
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) GitPush() *exec.Cmd {
	cmd := exec.Command("jj", "git", "push")
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) Diff(revision string) *exec.Cmd {
	cmd := exec.Command("jj", "diff", "-r", revision)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) Edit(revision string) *exec.Cmd {
	cmd := exec.Command("jj", "edit", "-r", revision)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) DiffEdit(revision string) *exec.Cmd {
	cmd := exec.Command("jj", "diffedit", "-r", revision)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) Split(revision string) *exec.Cmd {
	cmd := exec.Command("jj", "split", "-r", revision)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) Abandon(revision string) *exec.Cmd {
	cmd := exec.Command("jj", "abandon", "-r", revision)
	cmd.Dir = jj.Location
	return cmd
}

func (jj JJ) New(from string) *exec.Cmd {
	cmd := exec.Command("jj", "new", "-r", from)
	cmd.Dir = jj.Location
	return cmd
}
