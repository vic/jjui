package jj

import (
	"container/list"
	"fmt"
	"os/exec"
	"strings"
)

const TEMPLATE = `separate("\n", "__BEGIN__", change_id.shortest(1), change_id.short(8), coalesce(parents.map(|c| c.change_id().short(8)), "!!NONE"), current_working_copy, author.email(), coalesce(branches, "!!NONE"), coalesce(description, "!!NONE"), "__END__\n")`

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

func BuildCommitTree(commits []Commit) []Commit {
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

	visited := make(map[string]bool)
	stack := list.New()
	for i := len(commits) - 1; i >= 0; i-- {
		root := &commits[i]
		if _, ok := visited[root.ChangeId]; !ok {
			dfs(root, visited, stack, 0)
		}
	}
	commitsArray := make([]Commit, 0)
	// enumerate stack in reverse
	for i := stack.Back(); i != nil; i = i.Prev() {
		commitsArray = append(commitsArray, *i.Value.(*Commit))
	}
	return commitsArray
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
	return BuildCommitTree(commits)
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

func dfs(commit *Commit, visited map[string]bool, stack *list.List, level int) {
	commit.level = level
	visited[commit.ChangeId] = true
	stack.PushBack(commit)
	for i := len(commit.children) - 1; i >= 0; i-- {
		child := commit.children[i]
		if _, ok := visited[child.ChangeId]; !ok {
			dfs(child, visited, stack, level+i)
		}
	}
}

func RebaseCommand(from string, to string) error {
	cmd := exec.Command("jj", "rebase", "-r", from, "-d", to)
	err := cmd.Run()
	return err
}
