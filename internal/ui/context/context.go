package context

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/config"
)

type AppContext interface {
	KeyMap() config.KeyMappings[key.Binding]
	SelectedItem() SelectedItem
	SetSelectedItem(item SelectedItem)
	RunCommandImmediate(args []string) ([]byte, error)
	RunCommand(args []string, continuations ...tea.Cmd) tea.Cmd
	RunInteractiveCommand(args []string, continuation tea.Cmd) tea.Cmd
}
