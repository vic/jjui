package common

import (
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
	ShowOutputMsg      struct {
		Output string
		Err    error
	}
	SelectRevisionMsg string
)

type Operation int

const (
	None Operation = iota
	RebaseRevision
	RebaseBranch
)

func Close() tea.Msg {
	return CloseViewMsg{}
}

func SelectRevision(revision string) tea.Cmd {
	return func() tea.Msg {
		return SelectRevisionMsg(revision)
	}
}

func ShowOutput(output string, err error) tea.Cmd {
	return func() tea.Msg {
		return ShowOutputMsg{
			Output: output,
			Err:    err,
		}
	}
}

func ShowDescribe(commit *jj.Commit) tea.Cmd {
	return func() tea.Msg {
		return ShowDescribeViewMsg(commit)
	}
}

func RebaseCommand(from, to string, operation Operation) tea.Cmd {
	rebase := jj.RebaseCommand
	if operation == RebaseBranch {
		rebase = jj.RebaseBranchCommand
	}
	output, err := rebase(from, to)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")), SelectRevision(from))
}

func UpdateDescription(revision string, description string) tea.Cmd {
	output, err := jj.SetDescription(revision, description)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")))
}

func SetBookmark(revision string, bookmark string) tea.Cmd {
	output, err := jj.SetBookmark(revision, bookmark)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")))
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
