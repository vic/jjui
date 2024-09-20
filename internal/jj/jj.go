package jj

import (
	"fmt"
	"os/exec"
	"strings"
)

const TEMPLATE = `separate("\n", "__BEGIN__", change_id.shortest(1), change_id.short(8), coalesce(parents.map(|c| c.change_id().short(8)), "!!NONE"), current_working_copy, immutable, author.email(), coalesce(branches, "!!NONE"), coalesce(description, "!!NONE"), "__END__\n")`

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	Parents       []string
	IsWorkingCopy bool
	Author        string
	Branches      string
	Description   string
	Immutable     bool
	Index         int
}

func GetCommits(location string) []Commit {
	cmd := exec.Command("jj", "log", "--no-graph", "--template", TEMPLATE)
	cmd.Dir = location
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return parseLogOutput(string(output))
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
			commit.Index = i
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
	author := lines[6][indent:]
	if author != "!!NONE" {
		commit.Author = author
	}
	branches := lines[7][indent:]
	if branches != "!!NONE" {
		commit.Branches = branches
	}
	if len(lines) >= 9 {
		desc := lines[8][indent:]
		if desc != "!!NONE" {
			commit.Description = desc
		}
	}
	return commit
}

func RebaseCommand(from string, to string) error {
	cmd := exec.Command("jj", "rebase", "-r", from, "-d", to)
	err := cmd.Run()
	return err
}

func SetDescription(rev string, description string) error {
	cmd := exec.Command("jj", "describe", "-r", rev, "-m", description)
	err := cmd.Run()
	return err
}
