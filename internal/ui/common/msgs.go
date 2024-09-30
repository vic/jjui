package common

import (
	"os"
	"os/exec"

	"jjui/internal/dag"
	"jjui/internal/jj"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	ShowRevisionsViewMsg     struct{}
	CloseViewMsg             struct{}
	ShowDescribeOverlayMsg   *jj.Commit
	ShowAbandonOverlayMsg    string
	SelectRevisionMsg        string
	ShowDiffMsg              string
	UpdateRevisionsMsg       []dag.GraphRow
	UpdateBookmarksMsg       []jj.Bookmark
	UpdateDescriptionViewMsg struct {
		Description string
	}
	ShowOutputMsg struct {
		Output string
		Err    error
	}
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

func ShowDescribeOverlay(commit *jj.Commit) tea.Cmd {
	return func() tea.Msg {
		return ShowDescribeOverlayMsg(commit)
	}
}

func ShowAbandonOverlay(revision string) tea.Cmd {
	return func() tea.Msg {
		return ShowAbandonOverlayMsg(revision)
	}
}

func GitFetch() tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.GitFetch()
		return ShowOutputMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(f, FetchRevisions(os.Getenv("PWD")))
}

func GitPush() tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.GitPush()
		return ShowOutputMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(f, FetchRevisions(os.Getenv("PWD")))
}

func Rebase(from, to string, operation Operation) tea.Cmd {
	rebase := jj.RebaseCommand
	if operation == RebaseBranch {
		rebase = jj.RebaseBranchCommand
	}
	output, err := rebase(from, to)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")), SelectRevision(from))
}

func SetDescription(revision string, description string) tea.Cmd {
	output, err := jj.SetDescription(revision, description)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")))
}

func MoveBookmark(revision string, bookmark string) tea.Cmd {
	output, err := jj.MoveBookmark(revision, bookmark)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")))
}

func FetchRevisions(location string) tea.Cmd {
	return func() tea.Msg {
		commits, parents := jj.GetCommits(location)
		root := dag.Build(commits, parents)
		rows := dag.BuildGraphRows(root)
		return UpdateRevisionsMsg(rows)
	}
}

func FetchBookmarks(revision string) tea.Cmd {
	return func() tea.Msg {
		bookmarks, _ := jj.ListBookmark(revision)
		return UpdateBookmarksMsg(bookmarks)
	}
}

func GetDiff(revision string) tea.Cmd {
	return func() tea.Msg {
		output, _ := jj.Diff(revision)
		return ShowDiffMsg(output)
	}
}

func Edit(revision string) tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.Edit(revision)
		return ShowOutputMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(f, FetchRevisions(os.Getenv("PWD")), SelectRevision("@"))
}

func Split(revision string) tea.Cmd {
	return tea.Sequence(tea.ExecProcess(exec.Command("jj", "split", "-r", revision), nil), tea.ClearScreen, FetchRevisions(os.Getenv("PWD")))
}

func Abandon(revision string) tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.Abandon(revision)
		return ShowOutputMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(f, FetchRevisions(os.Getenv("PWD")))
}

func NewRevision(from string) tea.Cmd {
	output, err := jj.New(from)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")), SelectRevision("@"))
}
