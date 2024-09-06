package jj

import (
	"fmt"
	"os/exec"
	"strings"
)

// const TEMPLATE = `change_id.shortest(1) ++ "," ++ change_id.shortest(6) ++  "," ++ branch ++ "," ++ if(empty, "(empty)", description) ++ "\n"`
const TEMPLATE = `separate("\n", "__BEGIN__", change_id.shortest(1), change_id.shortest(6), commit_id, author, branches, coalesce(description, "empty"), "__END__\n")`

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
	fmt.Println(string(output))
	lines := strings.Split(string(output), "\n")
	start := -1
	commits := make([]Commit, 0)
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "__BEGIN__") {
			start = i
			continue
		}
		if strings.Contains(lines[i], "__END__") {
			fmt.Printf("start: %v, end: %v\n", start, i)
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
	commit.Author = lines[4][indent:]
	if len(lines) >= 7 {
		commit.Description = lines[6][indent:]
	}
	return commit
}
