package common

import (
	"fmt"
	"os"

	"jjui/internal/dag"
	"jjui/internal/jj"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	ShowRevisionsView struct{}
	CloseViewMsg      struct{}
)

type ShowDescribeViewMsg *jj.Commit

type UpdateDescriptionView struct {
	Description string
}

type (
	UpdateRevisionsMsg []dag.GraphRow
	UpdateBookmarksMsg []jj.Bookmark
)

func Close() tea.Msg {
	return CloseViewMsg{}
}

func ShowDescribe(commit *jj.Commit) tea.Cmd {
	return func() tea.Msg {
		return ShowDescribeViewMsg(commit)
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
		return UpdateRevisionsMsg(rows)
	}
}

func FetchBookmarks() tea.Msg {
	bookmarks, _ := jj.BookmarkList()
	return UpdateBookmarksMsg(bookmarks)
}
