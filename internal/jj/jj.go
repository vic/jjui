package jj

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type JJ struct {
	Location string
}

type Bookmark string

func (jj JJ) GetConfig(key string) ([]byte, error) {
	cmd := exec.Command("jj", "config", "get", key)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	output = bytes.Trim(output, "\n")
	return output, err
}

func (jj JJ) RebaseCommand(from string, to string) ([]byte, error) {
	cmd := exec.Command("jj", "rebase", "-r", from, "-d", to)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) RebaseBranchCommand(from string, to string) ([]byte, error) {
	cmd := exec.Command("jj", "rebase", "-b", from, "-d", to)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) SetDescription(rev string, description string) ([]byte, error) {
	cmd := exec.Command("jj", "describe", "-r", rev, "-m", description)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) ListBookmark(revision string) ([]Bookmark, error) {
	cmd := exec.Command("jj", "log", "-r", fmt.Sprintf("::%s- & bookmarks()", revision), "--template", "local_bookmarks.map(|x| x.name() ++ '\n')", "--no-graph")
	cmd.Dir = jj.Location
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var bookmarks []Bookmark
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		bookmarks = append(bookmarks, Bookmark(line))
	}
	return bookmarks, nil
}

func (jj JJ) SetBookmark(revision string, name string) ([]byte, error) {
	cmd := exec.Command("jj", "bookmark", "set", "-r", revision, name)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) MoveBookmark(revision string, bookmark string) ([]byte, error) {
	cmd := exec.Command("jj", "bookmark", "move", bookmark, "--to", revision)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) GitFetch() ([]byte, error) {
	cmd := exec.Command("jj", "git", "fetch")
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) GitPush() ([]byte, error) {
	cmd := exec.Command("jj", "git", "push")
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) Diff(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "diff", "-r", revision)
	cmd.Dir = jj.Location
	output, err := cmd.Output()
	return output, err
}

func (jj JJ) Edit(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "edit", "-r", revision)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) DiffEdit(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "diffedit", "-r", revision)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) Abandon(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "abandon", "-r", revision)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	return output, err
}

func (jj JJ) New(from string) ([]byte, error) {
	cmd := exec.Command("jj", "new", "-r", from)
	cmd.Dir = jj.Location
	output, err := cmd.Output()
	return output, err
}
