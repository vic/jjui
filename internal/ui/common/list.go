package common

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
)

type FilterableList struct {
	List          list.Model
	Items         []list.Item
	Filter        string
	KeyMap        config.KeyMappings[key.Binding]
	width         int
	height        int
	FilterMatches func(item list.Item, filter string) bool
	Title         string
	styles        styles
}

type styles struct {
	title    lipgloss.Style
	shortcut lipgloss.Style
	dimmed   lipgloss.Style
	selected lipgloss.Style
	matched  lipgloss.Style
	text     lipgloss.Style
	border   lipgloss.Style
}

type FilterMatchFunc func(list.Item, string) bool

func DefaultFilterMatch(item list.Item, filter string) bool {
	return true
}

func NewFilterableList(items []list.Item, width int, height int, keyMap config.KeyMappings[key.Binding]) FilterableList {
	styles := styles{
		title:    DefaultPalette.Get("menu title").Padding(0, 1, 0, 1),
		selected: DefaultPalette.Get("menu selected"),
		matched:  DefaultPalette.Get("menu matched"),
		dimmed:   DefaultPalette.Get("menu dimmed"),
		shortcut: DefaultPalette.Get("menu shortcut"),
		text:     DefaultPalette.Get("menu text"),
		border:   DefaultPalette.GetBorder("menu border", lipgloss.NormalBorder()),
	}

	delegate := ListItemDelegate{styles: styles}

	l := list.New(items, delegate, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(true)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	m := FilterableList{
		List:          l,
		Items:         items,
		KeyMap:        keyMap,
		FilterMatches: DefaultFilterMatch,
		styles:        styles,
	}
	m.SetWidth(width)
	m.SetHeight(height)

	l.Styles.NoItems = styles.dimmed
	l.Styles.PaginationStyle = styles.title.Width(10)
	l.Styles.ActivePaginationDot = styles.title
	l.Styles.InactivePaginationDot = styles.title
	l.FilterInput.PromptStyle = styles.matched
	l.FilterInput.Cursor.Style = styles.text

	return m
}

func (m *FilterableList) Width() int {
	return m.width
}

func (m *FilterableList) Height() int {
	return m.height
}

func (m *FilterableList) SetWidth(w int) {
	maxWidth, minWidth := 80, 40
	m.width = max(min(maxWidth, w), minWidth)
	m.List.SetWidth(m.width - 2)
}

func (m *FilterableList) SetHeight(h int) {
	maxHeight, minHeight := 30, 10
	m.height = max(min(maxHeight, h-2), minHeight)
	m.List.SetHeight(m.height - 2)
}

func (m *FilterableList) ShowShortcuts(show bool) {
	m.List.SetDelegate(ListItemDelegate{ShowShortcuts: show, styles: m.styles})
}

func (m *FilterableList) Filtered(filter string) tea.Cmd {
	m.Filter = filter
	if m.Filter == "" {
		m.List.SetDelegate(ListItemDelegate{ShowShortcuts: false, styles: m.styles})
		return m.List.SetItems(m.Items)
	}

	m.List.SetDelegate(ListItemDelegate{ShowShortcuts: true, styles: m.styles})
	var filtered []list.Item
	for _, i := range m.Items {
		if m.FilterMatches(i, m.Filter) {
			filtered = append(filtered, i)
		}
	}
	m.List.ResetSelected()
	return m.List.SetItems(filtered)
}

func (m *FilterableList) renderFilterView() string {
	filterStyle := m.styles.text.PaddingLeft(1)
	filterValueStyle := m.styles.matched

	filterView := lipgloss.JoinHorizontal(0, filterStyle.Render("Showing "), filterValueStyle.Render("all"))
	if m.Filter != "" {
		filterView = lipgloss.JoinHorizontal(0, filterStyle.Render("Showing only "), filterValueStyle.Render(m.Filter))
	}
	filterViewWidth := lipgloss.Width(filterView)
	paginationView := m.styles.text.AlignHorizontal(1).PaddingRight(1).Width(m.width - filterViewWidth).Render(fmt.Sprintf("%d/%d", m.List.Paginator.Page+1, m.List.Paginator.TotalPages))
	content := lipgloss.JoinHorizontal(0, filterView, paginationView)
	return m.styles.text.Width(m.width).Render(content)
}

func (m *FilterableList) renderHelpView(helpKeys []key.Binding) string {
	if m.List.SettingFilter() {
		return ""
	}

	bindings := make([]string, 0, len(helpKeys)+1)
	for _, k := range helpKeys {
		if renderedKey := m.renderKey(k); renderedKey != "" {
			bindings = append(bindings, renderedKey)
		}
	}

	if m.List.IsFiltered() {
		bindings = append(bindings, m.renderKey(m.KeyMap.Cancel))
	} else {
		bindings = append(bindings, m.renderKey(m.List.KeyMap.Filter))
	}

	return m.styles.text.PaddingLeft(1).Width(m.width).Render(lipgloss.JoinHorizontal(0, bindings...))
}

func (m *FilterableList) renderKey(k key.Binding) string {
	if !k.Enabled() {
		return ""
	}
	return lipgloss.JoinHorizontal(0, m.styles.shortcut.Render(k.Help().Key, ""), m.styles.dimmed.Render(k.Help().Desc, ""))
}

func (m *FilterableList) View(helpKeys []key.Binding) string {
	titleView := m.styles.text.Width(m.width).Render(m.styles.title.Render(m.Title))
	filterView := m.renderFilterView()
	listView := m.List.View()
	helpView := m.renderHelpView(helpKeys)
	content := lipgloss.JoinVertical(0, titleView, "", filterView, listView, helpView)
	content = lipgloss.Place(m.width, m.height, 0, 0, content)
	content = m.styles.text.Width(m.width).Height(m.height).Render(content)
	return m.styles.border.Render(content)
}
