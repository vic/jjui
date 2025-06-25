package squash

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	context context.AppContext
	From    jj.SelectedRevisions
	Current *jj.Commit
	keyMap  config.KeyMappings[key.Binding]
}

func (s *Operation) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, s.keyMap.Apply):
		return tea.Batch(common.Close, s.context.RunInteractiveCommand(jj.Squash(s.From, s.Current.ChangeId), common.Refresh))
	case key.Matches(msg, s.keyMap.Cancel):
		return common.Close
	}
	return nil
}

func (s *Operation) SetSelectedRevision(commit *jj.Commit) {
	s.Current = commit
}

func (s *Operation) Render(commit *jj.Commit, pos operations.RenderPosition) string {
	if pos != operations.RenderBeforeChangeId {
		return ""
	}

	isSelected := s.Current != nil && s.Current.GetChangeId() == commit.GetChangeId()
	if isSelected {
		return common.DefaultPalette.Drop.Render("<< into >> ")
	}
	sourceIds := s.From.GetIds()
	if slices.Contains(sourceIds, commit.ChangeId) {
		return common.DefaultPalette.EmptyPlaceholder.Render("<< from >> ")
	}
	return ""
}

func (s *Operation) Name() string {
	return "squash"
}

func (s *Operation) ShortHelp() []key.Binding {
	return []key.Binding{
		s.keyMap.Apply,
		s.keyMap.Cancel,
	}
}

func (s *Operation) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

func NewOperation(context context.AppContext, from jj.SelectedRevisions) *Operation {
	return &Operation{
		context: context,
		keyMap:  context.KeyMap(),
		From:    from,
	}
}
