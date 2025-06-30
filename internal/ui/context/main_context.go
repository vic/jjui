package context

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"log"
)

type SelectedItem interface {
	Equal(other SelectedItem) bool
}

type SelectedRevision struct {
	ChangeId string
}

func (s SelectedRevision) Equal(other SelectedItem) bool {
	if o, ok := other.(SelectedRevision); ok {
		return s.ChangeId == o.ChangeId
	}
	return false
}

type SelectedFile struct {
	ChangeId string
	File     string
}

func (s SelectedFile) Equal(other SelectedItem) bool {
	if o, ok := other.(SelectedFile); ok {
		return s.ChangeId == o.ChangeId && s.File == o.File
	}
	return false
}

type SelectedOperation struct {
	OperationId string
}

func (s SelectedOperation) Equal(other SelectedItem) bool {
	if o, ok := other.(SelectedOperation); ok {
		return s.OperationId == o.OperationId
	}
	return false
}

type MainContext struct {
	CommandRunner
	SelectedItem SelectedItem
	Location     string
	Config       *config.Config
	JJConfig     *jj.Config
}

func (a *MainContext) KeyMap() config.KeyMappings[key.Binding] {
	return a.Config.GetKeyMap()
}

func (a *MainContext) SetSelectedItem(item SelectedItem) tea.Cmd {
	if item == nil {
		return nil
	}
	if item.Equal(a.SelectedItem) {
		return nil
	}
	a.SelectedItem = item
	return common.SelectionChanged
}

func NewAppContext(location string) *MainContext {
	configuration, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	m := &MainContext{
		CommandRunner: &MainCommandRunner{
			Location: location,
		},
		Location: location,
		Config:   configuration,
	}

	m.JJConfig = &jj.Config{}
	if output, err := m.RunCommandImmediate(jj.ConfigListAll()); err == nil {
		if m.JJConfig, err = jj.DefaultConfig(output); err == nil {
			common.DefaultPalette.Update(m.JJConfig.Colors)
		}
	}
	return m
}
