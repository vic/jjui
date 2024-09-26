package jj

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	TEMPLATE             = `separate("\n", "__BEGIN__", change_id.shortest(1), change_id.short(8), coalesce(parents.map(|c| c.change_id().short(8)), "!!NONE"), current_working_copy, immutable, conflict, empty, author.email(), author.timestamp().ago(), coalesce(branches, "!!NONE"), coalesce(description, "!!NONE"), "__END__\n")`
	DESCENDANTS_TEMPLATE = `separate(" ", change_id.shortest(8), parents.map(|x| x.change_id().shortest(8))) ++ "\n"`
)

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	Parents       []string
	IsWorkingCopy bool
	Author        string
	Timestamp     string
	Branches      string
	Description   string
	Immutable     bool
	Conflict      bool
	Empty         bool
	Index         int
}

type Bookmark string

func GetCommits(location string) ([]Commit, map[string]string) {
	cmd := exec.Command("jj", "log", "--no-graph", "--template", TEMPLATE)
	cmd.Dir = location
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, nil
	}
	commits := parseLogOutput(string(output))
	parents := GetDescendants(commits[len(commits)-1].ChangeId)
	return commits, parents
}

func GetDescendants(root string) map[string]string {
	cmd := exec.Command("jj", "log", "--no-graph", "-r", root+"::", "--template", DESCENDANTS_TEMPLATE)
	cmd.Dir = os.Getenv("PWD")
	output, _ := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	parents := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		rev := parts[0]
		if len(parts) < 2 {
			parents[rev] = ""
			continue
		}
		parent := parts[1]
		parents[rev] = parent
	}
	return parents
}

func parseLogOutput(output string) []Commit {
	lines := strings.Split(output, "\n")
	start := -1
	commits := make([]Commit, 0)
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "__BEGIN__") {
			start = i
			continue
		}
		if strings.Contains(lines[i], "__END__") {
			commit := parseCommit(lines[start:i])
			commit.Index = len(commits)
			commits = append(commits, commit)
			start = -1
		}
	}
	return commits
}

func parseCommit(lines []string) Commit {
	indent := strings.Index(lines[0], "__BEGIN__")
	commit := Commit{}
	commit.ChangeIdShort = lines[1][indent:]
	commit.ChangeId = lines[2][indent:]
	parents := lines[3][indent:]
	if parents != "!!NONE" {
		commit.Parents = strings.Split(parents, " ")
	}
	commit.IsWorkingCopy = lines[4][indent:] == "true"
	commit.Immutable = lines[5][indent:] == "true"
	commit.Conflict = lines[6][indent:] == "true"
	commit.Empty = lines[7][indent:] == "true"
	author := lines[8][indent:]
	if author != "!!NONE" {
		commit.Author = author
	}
	commit.Timestamp = lines[9][indent:]
	branches := lines[10][indent:]
	if branches != "!!NONE" {
		commit.Branches = branches
	}
	if len(lines) >= 12 {
		desc := lines[11][indent:]
		if desc != "!!NONE" {
			commit.Description = desc
		}
	}
	return commit
}

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

func BookmarkList() ([]Bookmark, error) {
	cmd := exec.Command("jj", "bookmark", "list", "--template", "name ++ '\n'")
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

func SetBookmark(revision string, bookmark string) ([]byte, error) {
	cmd := exec.Command("jj", "bookmark", "set", bookmark, "-r", revision)
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

func New(from string) ([]byte, error) {
	cmd := exec.Command("jj", "new", "-r", from)
	output, err := cmd.Output()
	return output, err
}
