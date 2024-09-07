package jj

import (
	"fmt"
	"os/exec"
	"strings"
)

const TEMPLATE = `separate("\n", "__BEGIN__", change_id.shortest(1), change_id.short(), coalesce(parents.map(|c| c.change_id().short()), "!!NONE"), current_working_copy, author, coalesce(branches, "!!NONE"), coalesce(description, "!!NONE"), "__END__\n")`

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	Parent        string
	IsWorkingCopy bool
	Author        string
	Branches      string
	Description   string
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
			commits = append(commits, parseCommit(lines[start:i]))
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
	parent := lines[3][indent:]
	if parent != "!!NONE" {
		commit.Parent = parent
	}
	commit.IsWorkingCopy = lines[4][indent:] == "true"
	author := lines[5][indent:]
	if author != "!!NONE" {
		commit.Author = author
	}
	branches := lines[6][indent:]
	if branches != "!!NONE" {
		commit.Branches = branches
	}
	if len(lines) >= 8 {
		desc := lines[7][indent:]
		if desc != "!!NONE" {
			commit.Description = desc
		} else {
			commit.Description = "(empty)"
		}
	}
	return commit
}
