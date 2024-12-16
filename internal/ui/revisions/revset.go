package revisions

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/ui/common"
	"strings"
	"unicode"
)

var allowedFunctions = []string{
	"all()",
	"mine()",
	"empty()",
	"trunk()",
	"root()",
	"description(",
	"author(",
	"author_date(",
	"committer(",
	"committer_date(",
	"tags(",
	"files(",
	"bookmarks(",
	"conflicts()",
	"diff_contains(",
	"descendants(",
	"parents(",
	"ancestors(",
	"connected(",
	"git_head()",
	"git_refs()",
	"heads(",
	"fork_point(",
	"merges(",
	"remote_bookmarks(",
	"present(",
	"coalesce(",
	"working_copies()",
	"at_operation(",
	"builtin_immutable_heads()",
	"immutable()",
	"immutable_heads()",
	"mutable()",
	"tracked_remote_bookmarks(",
	"untracked_remote_bookmarks(",
	"visible_heads()",
	"latest(",
}

type EditRevSetMsg struct{}

type RevSetModel struct {
	Editing   bool
	Value     string
	textInput textinput.Model
}

var promptStyle = common.DefaultPalette.CommitShortStyle
var cursorStyle = common.DefaultPalette.Empty

func NewRevSet() RevSetModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = "revset: "
	ti.PromptStyle = promptStyle
	ti.Cursor.Style = cursorStyle
	ti.Focus()
	ti.CharLimit = 50
	ti.ShowSuggestions = true
	return RevSetModel{
		Editing:   false,
		Value:     "",
		textInput: ti,
	}
}

func (m RevSetModel) Init() tea.Cmd {
	return nil
}

func (m RevSetModel) Update(msg tea.Msg) (RevSetModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Editing = false
			return m, nil
		case tea.KeyEnter:
			m.Editing = false
			m.Value = m.textInput.Value()
			return m, tea.Batch(common.Close, common.UpdateRevSet(m.Value))
		}
	case EditRevSetMsg:
		m.Editing = true
		m.textInput.Focus()
		m.textInput.SetValue("")
		return m, textinput.Blink
	}

	value := m.textInput.Value()
	var suggestions []string
	lastIndex := strings.LastIndexFunc(value, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == '|' || r == '&' || r == '~'
	})

	if lastIndex == -1 && value == "" {
		suggestions = []string{"@ | mine()"}
	} else {
		lastFunctionName := value[lastIndex+1:]
		if !strings.HasSuffix(value, ")") && lastFunctionName != "" {
			for _, f := range allowedFunctions {
				if strings.HasPrefix(f, lastFunctionName) {
					rest := strings.TrimPrefix(f, lastFunctionName)
					suggestions = append(suggestions, value+rest)
				}
			}
		}
	}
	m.textInput.SetSuggestions(suggestions)

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m RevSetModel) View() string {
	if m.Editing {
		return m.textInput.View()
	}

	revset := "(default)"
	if m.Value != "" {
		revset = m.Value
	}

	return fmt.Sprintf("%s: %s", promptStyle.Render("revset"), cursorStyle.Render(revset))
}
