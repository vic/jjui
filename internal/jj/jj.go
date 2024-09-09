package jj

import (
	"container/list"
	"fmt"
	"os/exec"
	"strings"
)

const TEMPLATE = `separate("\n", "__BEGIN__", change_id.shortest(1), change_id.short(), coalesce(parents.map(|c| c.change_id().short()), "!!NONE"), current_working_copy, author, coalesce(branches, "!!NONE"), coalesce(description, "!!NONE"), "__END__\n")`

type Commit struct {
	ChangeIdShort string
	ChangeId      string
	Parents       []string
	IsWorkingCopy bool
	Author        string
	Branches      string
	Description   string
	children      []*Commit
	level         int
}

func (c Commit) Level() int {
	return c.level
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
	changeIdCommitMap := make(map[string]*Commit)
	for i, _ := range commits {
		commit := &commits[i]
		changeIdCommitMap[commit.ChangeId] = commit
	}
	for i, _ := range commits {
		commit := &commits[i]
		for _, parent := range commit.Parents {
			if parent, ok := changeIdCommitMap[parent]; ok {
				parent.children = append(parent.children, commit)
			}
		}
	}

	stack := dfsPushCommits(&commits[len(commits)-1])
	commitsArray := make([]Commit, 0, stack.Len())
	for e := stack.Front(); e != nil; e = e.Next() {
		commitsArray = append(commitsArray, *e.Value.(*Commit))
	}
	return commitsArray
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

func dfsPushCommits(root *Commit) *list.List {
	stack := list.New()
	dfs(root, stack, 0)
	return stack
}

func dfs(commit *Commit, stack *list.List, level int) {
	commit.level = level
	for i, child := range commit.children {
		dfs(child, stack, level+i)
	}
	stack.PushBack(commit)
}
