package details

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"
	"github.com/idursun/jjui/internal/ui/context"
)

type status uint8

var (
	Added    status = 0
	Deleted  status = 1
	Modified status = 2
	Renamed  status = 3
)

type item struct {
	status   status
	name     string
	fileName string
	selected bool
	conflict bool
}

func (f item) Title() string {
	status := "M"
	switch f.status {
	case Added:
		status = "A"
	case Deleted:
		status = "D"
	case Modified:
		status = "M"
	case Renamed:
		status = "R"
	}

	return fmt.Sprintf("%s %s", status, f.name)
}
func (f item) Description() string { return "" }
func (f item) FilterValue() string { return f.name }

type itemDelegate struct {
	selectedHint        string
	unselectedHint      string
	isVirtuallySelected bool
	styles              styles
}

func (i itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(item)
	if !ok {
		return
	}
	var style lipgloss.Style
	switch item.status {
	case Added:
		style = i.styles.Added
	case Deleted:
		style = i.styles.Deleted
	case Modified:
		style = i.styles.Modified
	case Renamed:
		style = i.styles.Renamed
	}

	if index == m.Index() {
		style = style.Bold(true).Background(i.styles.Selected.GetBackground())
	} else {
		style = style.Background(i.styles.Text.GetBackground())
	}

	title := item.Title()
	if item.selected {
		title = "✓" + title
	} else {
		title = " " + title
	}

	hint := ""
	if i.showHint() {
		hint = i.unselectedHint
		if item.selected || (i.isVirtuallySelected && index == m.Index()) {
			hint = i.selectedHint
			title = "✓" + item.Title()
		}
	}

	fmt.Fprint(w, style.PaddingRight(1).Render(title))
	if item.conflict {
		fmt.Fprint(w, i.styles.Conflict.Render("conflict "))
	}
	if hint != "" {
		fmt.Fprint(w, i.styles.Dimmed.Render(hint))
	}
}

func (i itemDelegate) Height() int                         { return 1 }
func (i itemDelegate) Spacing() int                        { return 0 }
func (i itemDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

func (i itemDelegate) showHint() bool {
	return i.selectedHint != "" || i.unselectedHint != ""
}

type styles struct {
	Added    lipgloss.Style
	Deleted  lipgloss.Style
	Modified lipgloss.Style
	Renamed  lipgloss.Style
	Selected lipgloss.Style
	Dimmed   lipgloss.Style
	Text     lipgloss.Style
	Conflict lipgloss.Style
}

type Model struct {
	revision     *jj.Commit
	files        list.Model
	height       int
	confirmation tea.Model
	context      *context.MainContext
	keyMap       config.KeyMappings[key.Binding]
	styles       styles
}

type updateCommitStatusMsg struct {
	summary       string
	selectedFiles []string
}

func New(context *context.MainContext, revision *jj.Commit) tea.Model {
	keyMap := config.Current.GetKeyMap()

	s := styles{
		Added:    common.DefaultPalette.Get("revisions details added"),
		Deleted:  common.DefaultPalette.Get("revisions details deleted"),
		Modified: common.DefaultPalette.Get("revisions details modified"),
		Renamed:  common.DefaultPalette.Get("revisions details renamed"),
		Selected: common.DefaultPalette.Get("revisions details selected"),
		Dimmed:   common.DefaultPalette.Get("revisions details dimmed"),
		Text:     common.DefaultPalette.Get("revisions details text"),
		Conflict: common.DefaultPalette.Get("revisions details conflict"),
	}

	l := list.New(nil, itemDelegate{styles: s}, 0, 0)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.KeyMap.CursorUp = keyMap.Up
	l.KeyMap.CursorDown = keyMap.Down
	l.Styles.NoItems = s.Dimmed
	return Model{
		revision: revision,
		files:    l,
		context:  context,
		keyMap:   keyMap,
		styles:   s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.load(m.revision.GetChangeId()), tea.WindowSize())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.confirmation != nil {
			model, cmd := m.confirmation.Update(msg)
			m.confirmation = model
			return m, cmd
		}
		switch {
		case key.Matches(msg, m.keyMap.Cancel), key.Matches(msg, m.keyMap.Details.Close):
			return m, common.Close
		case key.Matches(msg, m.keyMap.Details.Diff):
			selected, ok := m.files.SelectedItem().(item)
			if !ok {
				return m, nil
			}
			return m, func() tea.Msg {
				output, _ := m.context.RunCommandImmediate(jj.Diff(m.revision.GetChangeId(), selected.fileName))
				return common.ShowDiffMsg(output)
			}
		case key.Matches(msg, m.keyMap.Details.Edit):
			editRevision := func() tea.Msg { return common.EditRevision(m.revision) }
			selectedFiles, _ := m.getSelectedFiles()
			editFiles := func() tea.Msg {
				line := []string{config.GetDefaultEditor()}
				line = append(line, selectedFiles...)
				return common.ExecMsg{
					Line: strings.Join(line, " "),
					Mode: common.ExecShell,
				}
			}
			return m, tea.Batch(editRevision, editFiles)
		case key.Matches(msg, m.keyMap.Details.Split):
			selectedFiles, isVirtuallySelected := m.getSelectedFiles()
			m.files.SetDelegate(itemDelegate{
				isVirtuallySelected: isVirtuallySelected,
				selectedHint:        "stays as is",
				unselectedHint:      "moves to the new revision",
				styles:              m.styles,
			})
			model := confirmation.New(
				[]string{"Are you sure you want to split the selected files?"},
				confirmation.WithStylePrefix("revisions"),
				confirmation.WithOption("Yes", tea.Batch(m.context.RunInteractiveCommand(jj.Split(m.revision.GetChangeId(), selectedFiles), common.Refresh), common.Close), key.NewBinding(key.WithKeys("y"))),
				confirmation.WithOption("No", confirmation.Close, key.NewBinding(key.WithKeys("n", "esc"))),
			)
			m.confirmation = &model
			return m, m.confirmation.Init()
		case key.Matches(msg, m.keyMap.Details.Restore):
			selectedFiles, isVirtuallySelected := m.getSelectedFiles()
			m.files.SetDelegate(itemDelegate{
				isVirtuallySelected: isVirtuallySelected,
				selectedHint:        "gets restored",
				unselectedHint:      "stays as is",
				styles:              m.styles,
			})
			model := confirmation.New(
				[]string{"Are you sure you want to restore the selected files?"},
				confirmation.WithStylePrefix("revisions"),
				confirmation.WithOption("Yes", m.context.RunCommand(jj.Restore(m.revision.GetChangeId(), selectedFiles), common.Refresh, confirmation.Close), key.NewBinding(key.WithKeys("y"))),
				confirmation.WithOption("No", confirmation.Close, key.NewBinding(key.WithKeys("n", "esc"))),
			)
			m.confirmation = &model
			return m, m.confirmation.Init()
		case key.Matches(msg, m.keyMap.Details.Absorb):
			selectedFiles, isVirtuallySelected := m.getSelectedFiles()
			m.files.SetDelegate(itemDelegate{
				isVirtuallySelected: isVirtuallySelected,
				selectedHint:        "might get absorbed into parents",
				unselectedHint:      "stays as is",
				styles:              m.styles,
			})
			model := confirmation.New(
				[]string{"Are you sure you want to absorb changes from the selected files?"},
				confirmation.WithStylePrefix("revisions"),
				confirmation.WithOption("Yes", m.context.RunCommand(jj.Absorb(m.revision.GetChangeId(), selectedFiles...), common.Refresh, confirmation.Close), key.NewBinding(key.WithKeys("y"))),
				confirmation.WithOption("No", confirmation.Close, key.NewBinding(key.WithKeys("n", "esc"))),
			)
			m.confirmation = &model
			return m, m.confirmation.Init()
		case key.Matches(msg, m.keyMap.Details.ToggleSelect):
			if oldItem, ok := m.files.SelectedItem().(item); ok {
				oldItem.selected = !oldItem.selected
				oldIndex := m.files.Index()
				m.files.CursorDown()

				curItem := m.files.SelectedItem().(item)
				return m, tea.Batch(m.files.SetItem(oldIndex, oldItem), m.context.SetSelectedItem(context.SelectedFile{
					ChangeId: m.revision.GetChangeId(),
					CommitId: m.revision.CommitId,
					File:     curItem.fileName,
				}))
			}
			return m, nil
		case key.Matches(msg, m.keyMap.Details.RevisionsChangingFile):
			if item, ok := m.files.SelectedItem().(item); ok {
				return m, tea.Batch(common.Close, common.UpdateRevSet(fmt.Sprintf("files(\"%s\")", item.fileName)))
			}
		default:
			if len(m.files.Items()) > 0 {
				var cmd tea.Cmd
				m.files, cmd = m.files.Update(msg)
				curItem := m.files.SelectedItem().(item)
				return m, tea.Batch(cmd, m.context.SetSelectedItem(context.SelectedFile{
					ChangeId: m.revision.GetChangeId(),
					CommitId: m.revision.CommitId,
					File:     curItem.fileName,
				}))
			}
		}
	case confirmation.CloseMsg:
		m.confirmation = nil
		m.files.SetDelegate(itemDelegate{styles: m.styles})
		return m, nil
	case common.RefreshMsg:
		return m, m.load(m.revision.GetChangeId())
	case updateCommitStatusMsg:
		items := m.createListItems(msg.summary, msg.selectedFiles)
		var selectionChangedCmd tea.Cmd
		if len(items) > 0 {
			selectionChangedCmd = m.context.SetSelectedItem(context.SelectedFile{
				ChangeId: m.revision.GetChangeId(),
				CommitId: m.revision.CommitId,
				File:     items[0].(item).fileName,
			})
		}
		return m, tea.Batch(selectionChangedCmd, m.files.SetItems(items))
	case tea.WindowSizeMsg:
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) createListItems(content string, selectedFiles []string) []list.Item {
	items := make([]list.Item, 0)
	scanner := bufio.NewScanner(strings.NewReader(content))
	var conflicts []bool
	if scanner.Scan() {
		conflictsLine := strings.Split(scanner.Text(), " ")
		for _, c := range conflictsLine {
			conflicts = append(conflicts, c == "true")
		}
	} else {
		return items
	}

	index := 0
	for scanner.Scan() {
		file := strings.TrimSpace(scanner.Text())
		if file == "" {
			continue
		}
		var status status
		switch file[0] {
		case 'A':
			status = Added
		case 'D':
			status = Deleted
		case 'M':
			status = Modified
		case 'R':
			status = Renamed
		}
		fileName := file[2:]

		actualFileName := fileName
		if status == Renamed && strings.Contains(actualFileName, "{") {
			for strings.Contains(actualFileName, "{") {
				start := strings.Index(actualFileName, "{")
				end := strings.Index(actualFileName, "}")
				if end == -1 {
					break
				}
				replacement := actualFileName[start+1 : end]
				parts := strings.Split(replacement, " => ")
				replacement = parts[1]
				actualFileName = path.Clean(actualFileName[:start] + replacement + actualFileName[end+1:])
			}
		}
		items = append(items, item{
			status:   status,
			name:     fileName,
			fileName: actualFileName,
			selected: slices.ContainsFunc(selectedFiles, func(s string) bool { return s == actualFileName }),
			conflict: conflicts[index],
		})
		index++
	}
	return items
}

func (m Model) getSelectedFiles() ([]string, bool) {
	selectedFiles := make([]string, 0)
	if len(m.files.Items()) == 0 {
		return selectedFiles, false
	}

	isVirtuallySelected := false
	for _, f := range m.files.Items() {
		if f.(item).selected {
			selectedFiles = append(selectedFiles, f.(item).fileName)
			isVirtuallySelected = false
		}
	}
	if len(selectedFiles) == 0 {
		selectedFiles = append(selectedFiles, m.files.SelectedItem().(item).fileName)
		return selectedFiles, true
	}
	return selectedFiles, isVirtuallySelected
}

func (m Model) View() string {
	confirmationView := ""
	ch := 0
	if m.confirmation != nil {
		confirmationView = m.confirmation.View()
		ch = lipgloss.Height(confirmationView)
	}
	m.files.SetHeight(min(m.height-5-ch, len(m.files.Items())))
	filesView := m.files.View()

	view := lipgloss.JoinVertical(lipgloss.Top, filesView, confirmationView)
	// We are trimming spaces from each line to prevent visual artefacts
	// Empty lines use the default background colour, and it looks bad if the user has a custom background colour
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(view))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lines = append(lines, line)
	}
	view = strings.Join(lines, "\n")
	w, h := lipgloss.Size(view)
	return lipgloss.Place(w, h, 0, 0, view, lipgloss.WithWhitespaceBackground(m.styles.Text.GetBackground()))
}

func (m Model) load(revision string) tea.Cmd {
	output, err := m.context.RunCommandImmediate(jj.Snapshot())
	if err == nil {
		output, err = m.context.RunCommandImmediate(jj.Status(revision))
		if err == nil {
			return func() tea.Msg {
				summary := string(output)
				selectedFiles, isVirtuallySelected := m.getSelectedFiles()
				if isVirtuallySelected {
					selectedFiles = []string{}
				}
				return updateCommitStatusMsg{summary, selectedFiles}
			}
		}
	}
	return func() tea.Msg {
		return common.CommandCompletedMsg{
			Output: string(output),
			Err:    err,
		}
	}
}
