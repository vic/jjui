package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/jj"
)

type (
	CloseViewMsg             struct{}
	RefreshMsg               struct{ SelectedRevision string }
	SelectRevisionMsg        string
	ShowDiffMsg              string
	UpdateRevSetMsg          string
	UpdateRevisionsMsg       []jj.GraphLine
	UpdateRevisionsFailedMsg error
	UpdateBookmarksMsg       []jj.Bookmark
	CommandRunningMsg        string
	AbandonMsg               string
	SetDescriptionMsg        struct {
		Revision    string
		Description string
	}
	SetBookmarkMsg struct {
		Revision string
		Bookmark string
	}
	MoveBookmarkMsg struct {
		Revision string
		Bookmark string
	}
	CommandCompletedMsg struct {
		Output string
		Err    error
	}
)

type Operation int

const (
	None Operation = iota
	RebaseRevisionOperation
	RebaseBranchOperation
	EditDescriptionOperation
	SetBookmarkOperation
)

type Status int

const (
	Loading Status = iota
	Ready
	Error
)

func Close() tea.Msg {
	return CloseViewMsg{}
}

func Refresh(selectedRevision string) tea.Cmd {
	return func() tea.Msg {
		return RefreshMsg{SelectedRevision: selectedRevision}
	}
}

func UpdateRevSet(revset string) tea.Cmd {
	return func() tea.Msg {
		return UpdateRevSetMsg(revset)
	}
}

func SelectRevision(revision string) tea.Cmd {
	return func() tea.Msg {
		return SelectRevisionMsg(revision)
	}
}

func CommandRunning(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandRunningMsg(command)
	}
}

func ShowOutput(output string, err error) tea.Cmd {
	return func() tea.Msg {
		return CommandCompletedMsg{
			Output: output,
			Err:    err,
		}
	}
}
