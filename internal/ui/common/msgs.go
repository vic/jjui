package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/ui/operations"
	"strings"
)

type (
	CloseViewMsg             struct{}
	ToggleHelpMsg            struct{}
	RefreshMsg               struct{ SelectedRevision string }
	SelectRevisionMsg        string
	SetOperationMsg          struct{ Operation operations.Operation }
	ShowDiffMsg              string
	UpdateRevisionsFailedMsg error
	UpdateBookmarksMsg       struct {
		Bookmarks []string
		Revision  string
	}
	CommandRunningMsg   string
	CommandCompletedMsg struct {
		Output string
		Err    error
	}
	SelectionChangedMsg struct{}
)

type State int

const (
	Loading State = iota
	Ready
	Error
)

func Close() tea.Msg {
	return CloseViewMsg{}
}

func SelectionChanged() tea.Msg {
	return SelectionChangedMsg{}
}

func Refresh(selectedRevision string) tea.Cmd {
	return func() tea.Msg {
		return RefreshMsg{SelectedRevision: selectedRevision}
	}
}

func ToggleHelp() tea.Msg {
	return ToggleHelpMsg{}
}

func SelectRevision(revision string) tea.Cmd {
	return func() tea.Msg {
		return SelectRevisionMsg(revision)
	}
}

func SetOperation(op operations.Operation) tea.Cmd {
	return func() tea.Msg {
		return SetOperationMsg{Operation: op}
	}
}

func CommandRunning(args []string) tea.Cmd {
	return func() tea.Msg {
		command := "jj " + strings.Join(args, " ")
		return CommandRunningMsg(command)
	}
}
