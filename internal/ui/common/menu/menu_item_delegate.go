package menu

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
)

type MenuItem interface {
	list.DefaultItem
	ShortCut() string
}

type MenuItemDelegate struct {
	ShowShortcuts bool
	styles        styles
}

func (l MenuItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title    string
		desc     string
		shortcut string
	)
	if item, ok := item.(MenuItem); ok {
		title = item.Title()
		desc = item.Description()
		shortcut = item.ShortCut()
	} else {
		return
	}
	if m.Width() <= 0 {
		// short-circuit
		return
	}

	if !l.ShowShortcuts {
		shortcut = ""
	}

	titleWidth := m.Width()
	if shortcut != "" {
		titleWidth -= lipgloss.Width(shortcut) + 1
	}

	if len(title) > titleWidth {
		title = title[:titleWidth-1] + "…"
	}

	if len(desc) > m.Width() {
		desc = desc[:m.Width()-1] + "…"
	}

	var (
		titleStyle    = l.styles.text
		descStyle     = l.styles.dimmed
		shortcutStyle = l.styles.shortcut
		rowStyle      = l.styles.text
	)

	if index == m.Index() {
		titleStyle = l.styles.selected
		descStyle = l.styles.selected
		shortcutStyle = shortcutStyle.Background(l.styles.selected.GetBackground())
		rowStyle = l.styles.selected
	}

	rowStyle = rowStyle.PaddingLeft(1).PaddingRight(1).Width(m.Width() + 2)

	if shortcut != "" {
		_, _ = fmt.Fprintln(w, rowStyle.Render(lipgloss.JoinHorizontal(0, shortcutStyle.Render(shortcut), titleStyle.Render("", title))))
	} else {
		_, _ = fmt.Fprintln(w, rowStyle.Render(titleStyle.Render(title)))
	}

	descStyle = descStyle.PaddingLeft(1).PaddingRight(1).Width(m.Width() + 2)
	_, _ = fmt.Fprintf(w, descStyle.Render(desc))
}

func (l MenuItemDelegate) Height() int {
	return 2
}

func (l MenuItemDelegate) Spacing() int {
	return 1
}

func (l MenuItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
