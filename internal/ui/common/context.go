package common

import "github.com/idursun/jjui/internal/jj"

type SelectedItem interface{}

type SelectedRevision struct {
	ChangeId string
}

type SelectedFile struct {
	ChangeId string
	File     string
}

type AppContext struct {
	JJ           jj.Commands
	UICommands   UICommands
	SelectedItem SelectedItem
}

func (a *AppContext) SetSelectedItem(item SelectedItem) {
	a.SelectedItem = item
}

func NewAppContext(jj jj.Commands) *AppContext {
	return &AppContext{
		JJ:         jj,
		UICommands: NewUICommands(jj),
	}
}
