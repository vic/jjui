package context

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"strings"
)

type CustomRevsetCommand struct {
	CustomCommandBase
	Revset string `toml:"revset"`
}

func (c CustomRevsetCommand) Description(ctx *MainContext) string {
	if item, ok := ctx.SelectedItem.(SelectedRevision); ok {
		rendered := strings.ReplaceAll(c.Revset, jj.ChangeIdPlaceholder, item.ChangeId)
		return fmt.Sprintf("change revset to %s", rendered)
	}
	return ""
}

func (c CustomRevsetCommand) IsApplicableTo(item SelectedItem) bool {
	_, ok := item.(SelectedRevision)
	return ok
}

func (c CustomRevsetCommand) Prepare(ctx *MainContext) tea.Cmd {
	if item, ok := ctx.SelectedItem.(SelectedRevision); ok {
		rendered := strings.ReplaceAll(c.Revset, jj.ChangeIdPlaceholder, item.ChangeId)
		return common.UpdateRevSet(rendered)
	}
	return nil
}
