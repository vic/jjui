package jj

import (
	"fmt"
	"os/exec"
	"strings"
)

const TEMPLATE = `separate("\n", "__BEGIN__", change_id.shortest(1), change_id.shortest(6), commit_id, author, coalesce(branches, "!!NONE"), coalesce(description, "!!NONE"), "__END__\n")`

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	CommitId      string
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
	commit.CommitId = lines[3][indent:]
	author := lines[4][indent:]
	if author != "!!NONE" {
		commit.Author = author
	}
	branches := lines[5][indent:]
	if branches != "!!NONE" {
		commit.Branches = branches
	}
	if len(lines) >= 7 {
		desc := lines[6][indent:]
		if desc != "!!NONE" {
			commit.Description = desc
		} else {
			commit.Description = "(empty)"
		}
	}
	return commit
}
