package common

import (
	"fmt"
	"os"

	"jjui/internal/dag"
	"jjui/internal/jj"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	ShowRevisions struct{}
	CloseModal    struct{}
)

type ShowDescribeModal *jj.Commit

type UpdateDescriptionMessage struct {
	Description string
}

type UpdateRevisions []dag.GraphRow

func Close() tea.Msg {
	return CloseModal{}
}

func DoShowDescribe(commit *jj.Commit) tea.Cmd {
	return func() tea.Msg {
		return ShowDescribeModal(commit)
	}
}

func RebaseCommand(from, to string) tea.Cmd {
	if err := jj.RebaseCommand(from, to); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return FetchRevisions(os.Getenv("PWD"))
}

func UpdateDescription(revision string, description string) tea.Cmd {
	if err := jj.SetDescription(revision, description); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return FetchRevisions(os.Getenv("PWD"))
}

func FetchRevisions(location string) tea.Cmd {
	return func() tea.Msg {
		commits := jj.GetCommits(location)
		root := dag.Build(commits)
		rows := dag.BuildGraphRows(root)
		return UpdateRevisions(rows)
	}
}
