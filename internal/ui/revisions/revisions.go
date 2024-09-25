package revisions

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"jjui/internal/dag"
	"jjui/internal/jj"
	"jjui/internal/ui/bookmark"
	"jjui/internal/ui/common"
	"jjui/internal/ui/describe"
	"jjui/internal/ui/diff"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	normalMode mode = iota
	dropMode
)

type Model struct {
	rows       []dag.GraphRow
	op         common.Operation
	draggedRow int
	cursor     int
	width      int
	height     int
	help       help.Model
	describe   tea.Model
	bookmarks  tea.Model
	diff       tea.Model
	output     string
	keymap     keymap
}

type keymap struct {
	current  rune
	bindings map[rune]interface{}
	up       key.Binding
	down     key.Binding
	cancel   key.Binding
	apply    key.Binding
}

type baseLayer struct {
	edit         key.Binding
	rebaseMode   key.Binding
	bookmarkMode key.Binding
	gitMode      key.Binding
	description  key.Binding
	diff         key.Binding
	new          key.Binding
	quit         key.Binding
}

type rebaseLayer struct {
	revision key.Binding
	branch   key.Binding
}

type bookmarkLayer struct {
	set    key.Binding
	delete key.Binding
}

type gitLayer struct {
	fetch key.Binding
	push  key.Binding
}

func newKeyMap() keymap {
	bindings := make(map[rune]interface{})
	bindings[' '] = baseLayer{
		edit:         key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		rebaseMode:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rebase")),
		bookmarkMode: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "bookmark")),
		gitMode:      key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "git")),
		description:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "description")),
		diff:         key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "show diff")),
		new:          key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	bindings['r'] = rebaseLayer{
		revision: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rebase revision")),
		branch:   key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "rebase branch")),
	}

	bindings['b'] = bookmarkLayer{
		set:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "bookmark set")),
		delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "bookmark delete")),
	}

	bindings['g'] = gitLayer{
		fetch: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "git push")),
		push:  key.NewBinding(key.WithKeys("f"), key.WithHelp("d", "git fetch")),
	}

	return keymap{
		current:  ' ',
		bindings: bindings,
		up:       key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		down:     key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		apply:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
		cancel:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (k *keymap) gitMode() {
	k.current = 'g'
}

func (k *keymap) rebaseMode() {
	k.current = 'r'
}

func (k *keymap) bookmarkMode() {
	k.current = 'b'
}

func (k *keymap) resetMode() {
	k.current = ' '
}

func (k *keymap) ShortHelp() []key.Binding {
	switch b := k.bindings[k.current].(type) {
	case baseLayer:
		return []key.Binding{k.up, k.down, b.description, b.new, b.edit, b.rebaseMode, b.gitMode, b.quit}
	case rebaseLayer:
		return []key.Binding{b.revision, b.branch}
	case gitLayer:
		return []key.Binding{b.push, b.fetch}
	case bookmarkLayer:
		return []key.Binding{b.set, b.delete}
	default:
		return []key.Binding{}
	}
}

func (k *keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func (m Model) selectedRevision() *jj.Commit {
	return m.rows[m.cursor].Commit
}

func (m Model) Init() tea.Cmd {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return tea.Sequence(tea.SetWindowTitle("jjui"), common.FetchRevisions(dir), common.SelectRevision("@"))
}

func (m Model) handleBaseKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(baseLayer)
	switch {
	case key.Matches(msg, layer.quit):
		return m, tea.Quit
	case key.Matches(msg, layer.gitMode):
		m.keymap.gitMode()
		return m, nil
	case key.Matches(msg, layer.rebaseMode):
		m.keymap.rebaseMode()
		return m, nil
	case key.Matches(msg, layer.bookmarkMode):
		m.keymap.bookmarkMode()
		return m, nil
	case key.Matches(msg, layer.description):
		return m, common.ShowDescribe(m.selectedRevision())
	case key.Matches(msg, layer.diff):
		return m, common.GetDiff(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.new):
		return m, common.NewRevision(m.selectedRevision().ChangeId)
	case key.Matches(msg, m.keymap.down):
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case key.Matches(msg, m.keymap.up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.keymap.cancel):
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleRebaseKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(rebaseLayer)
	switch {
	case key.Matches(msg, layer.revision):
		// rebase revision
		m.op = common.RebaseRevision
		m.draggedRow = m.cursor
	case key.Matches(msg, layer.branch):
		m.keymap.current = 'm'
		m.op = common.RebaseBranch
		m.draggedRow = m.cursor
	case key.Matches(msg, m.keymap.down):
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case key.Matches(msg, m.keymap.up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.keymap.apply):
		m.op = common.None
		fromRevision := m.rows[m.draggedRow].Commit.ChangeIdShort
		toRevision := m.rows[m.cursor].Commit.ChangeIdShort
		m.draggedRow = -1
		m.keymap.current = ' '
		return m, common.RebaseCommand(fromRevision, toRevision, m.op)
	case key.Matches(msg, m.keymap.cancel):
		m.keymap.resetMode()
		m.op = common.None
	}
	return m, nil
}

func (m Model) handleBookmarkKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(bookmarkLayer)
	switch {
	case key.Matches(msg, layer.set):
		m.keymap.resetMode()
		return m, common.FetchBookmarks
	case key.Matches(msg, m.keymap.cancel):
		m.keymap.resetMode()
	}
	return m, nil
}

func (m Model) handleGitKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(gitLayer)
	switch {
	case key.Matches(msg, layer.fetch):
		m.keymap.resetMode()
		return m, common.GitFetch()
	case key.Matches(msg, layer.push):
		m.keymap.resetMode()
		return m, common.GitPush()
	case key.Matches(msg, m.keymap.cancel):
		m.keymap.resetMode()
	}
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(common.CloseViewMsg); ok {
		m.describe = nil
		m.bookmarks = nil
		m.diff = nil
		return m, nil
	}

	var cmd tea.Cmd
	if m.describe != nil {
		m.describe, cmd = m.describe.Update(msg)
		return m, cmd
	}
	if m.bookmarks != nil {
		m.bookmarks, cmd = m.bookmarks.Update(msg)
		return m, cmd
	}
	if m.diff != nil {
		m.diff, cmd = m.diff.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.keymap.current {
		case ' ':
			return m.handleBaseKeys(msg)
		case 'r':
			return m.handleRebaseKeys(msg)
		case 'b':
			return m.handleBookmarkKeys(msg)
		case 'g':
			return m.handleGitKeys(msg)
		}

	case common.SelectRevisionMsg:
		r := string(msg)
		idx := slices.IndexFunc(m.rows, func(row dag.GraphRow) bool {
			if r == "@" {
				return row.Commit.IsWorkingCopy
			}
			return row.Commit.ChangeIdShort == r
		})
		if idx != -1 {
			m.cursor = idx
		} else {
			m.cursor = 0
		}
	case common.UpdateRevisionsMsg:
		m.rows = msg
	case common.UpdateBookmarksMsg:
		m.bookmarks = bookmark.New(m.selectedRevision().ChangeId, msg, m.width)
		return m, m.bookmarks.Init()
	case common.ShowDescribeViewMsg:
		m.describe = describe.New(msg.ChangeId, msg.Description, m.width)
		return m, m.describe.Init()
	case common.ShowDiffMsg:
		m.diff = diff.New(string(msg), m.width, m.height)
		return m, m.diff.Init()
	case common.ShowOutputMsg:
		m.output = msg.Output
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, cmd
}

func (m Model) View() string {
	if m.height == 0 {
		return "loading"
	}

	if m.diff != nil {
		return m.diff.View()
	}

	var b strings.Builder
	b.WriteString(m.help.View(&m.keymap))
	b.WriteString("\n")
	if m.op == common.RebaseBranch || m.op == common.RebaseRevision {
		command := "-r"
		if m.op == common.RebaseBranch {
			command = "-b"
		}
		b.WriteString("jj rebase " + command + " " + m.rows[m.draggedRow].Commit.ChangeIdShort + " -d " + m.rows[m.cursor].Commit.ChangeIdShort + "\n")
	}

	if m.output != "" {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%v\n", m.output))
	}

	footer := b.String()
	footerHeight := lipgloss.Height(footer)
	if m.bookmarks != nil {
		footerHeight += 4
	}
	if m.describe != nil {
		footerHeight += 2
	}
	itemsPerPage := (m.height - footerHeight - 6) / 2
	if itemsPerPage == 0 {
		return b.String()
	}
	currentPage := m.cursor / itemsPerPage
	viewStart := currentPage * itemsPerPage
	viewEnd := (currentPage + 1) * itemsPerPage

	var items strings.Builder
	for i := 0; i < len(m.rows); i++ {
		if i < viewStart {
			continue
		}
		if i > viewEnd {
			continue
		}
		row := &m.rows[i]
		switch m.op {
		case common.RebaseRevision, common.RebaseBranch:
			if i == m.cursor {
				indent := strings.Repeat("│ ", row.Level)
				items.WriteString(indent)
				items.WriteString(common.DropStyle.Render("<< here >>"))
				items.WriteString("\n")
			}
			DefaultRenderer(&items, row, common.DefaultPalette, i == m.draggedRow)
		case common.None:
			DefaultRenderer(&items, row, common.DefaultPalette, i == m.cursor)
			if m.describe != nil && m.cursor == i {
				items.WriteString(m.describe.View())
				items.WriteString("\n")
			}
			if m.bookmarks != nil && m.cursor == i {
				items.WriteString(m.bookmarks.View())
				items.WriteString("\n")
			}
		}
	}
	items.WriteString("\n")
	items.WriteString(footer)
	return items.String()
}

func New(rows []dag.GraphRow) tea.Model {
	h := help.New()
	h.Styles.ShortKey = common.DefaultPalette.CommitShortStyle
	h.Styles.ShortDesc = common.DefaultPalette.CommitIdRestStyle
	return Model{
		rows:       rows,
		draggedRow: -1,
		op:         common.None,
		cursor:     0,
		width:      20,
		keymap:     newKeyMap(),
		describe:   nil,
		help:       h,
	}
}
