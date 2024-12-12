package jj

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	TEMPLATE = `separate(";", change_id.shortest(1), change_id.shortest(8), separate(",", parents.map(|x| x.change_id().shortest(1))), separate(",", coalesce(bookmarks, ".")), current_working_copy, immutable, conflict, empty, author.email(), author.timestamp().ago(), description.first_line())`
)

type Bookmark string

func RebaseCommand(from string, to string) ([]byte, error) {
	cmd := exec.Command("jj", "rebase", "-r", from, "-d", to)
	output, err := cmd.CombinedOutput()
	return output, err
}

func RebaseBranchCommand(from string, to string) ([]byte, error) {
	cmd := exec.Command("jj", "rebase", "-b", from, "-d", to)
	output, err := cmd.CombinedOutput()
	return output, err
}

func SetDescription(rev string, description string) ([]byte, error) {
	cmd := exec.Command("jj", "describe", "-r", rev, "-m", description)
	output, err := cmd.CombinedOutput()
	return output, err
}

func ListBookmark(revision string) ([]Bookmark, error) {
	cmd := exec.Command("jj", "log", "-r", fmt.Sprintf("::%s- & bookmarks()", revision), "--template", "local_bookmarks.map(|x| x.name() ++ '\n')", "--no-graph")
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

func MoveBookmark(revision string, bookmark string) ([]byte, error) {
	cmd := exec.Command("jj", "bookmark", "move", bookmark, "--to", revision)
	output, err := cmd.CombinedOutput()
	return output, err
}

func GitFetch() ([]byte, error) {
	cmd := exec.Command("jj", "git", "fetch")
	output, err := cmd.CombinedOutput()
	return output, err
}

func GitPush() ([]byte, error) {
	cmd := exec.Command("jj", "git", "push")
	output, err := cmd.CombinedOutput()
	return output, err
}

func Diff(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "diff", "-r", revision)
	output, err := cmd.Output()
	return output, err
}

func Edit(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "edit", "-r", revision)
	output, err := cmd.CombinedOutput()
	return output, err
}

func DiffEdit(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "diffedit", "-r", revision)
	output, err := cmd.CombinedOutput()
	return output, err
}

func Abandon(revision string) ([]byte, error) {
	cmd := exec.Command("jj", "abandon", "-r", revision)
	output, err := cmd.CombinedOutput()
	return output, err
}

func New(from string) ([]byte, error) {
	cmd := exec.Command("jj", "new", "-r", from)
	output, err := cmd.Output()
	return output, err
}
