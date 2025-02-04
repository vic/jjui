package revisions

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	current  rune
	bindings map[rune]interface{}
	up       key.Binding
	down     key.Binding
	details  key.Binding
	cancel   key.Binding
	apply    key.Binding
}

type baseLayer struct {
	abandon      key.Binding
	edit         key.Binding
	diffedit     key.Binding
	split        key.Binding
	rebaseMode   key.Binding
	squashMode   key.Binding
	bookmarkMode key.Binding
	gitMode      key.Binding
	description  key.Binding
	diff         key.Binding
	new          key.Binding
	revset       key.Binding
	refresh      key.Binding
	undo         key.Binding
	quit         key.Binding
}

type rebaseLayer struct {
	revision key.Binding
	branch   key.Binding
}

type squashLayer struct {
	apply key.Binding
}

type bookmarkLayer struct {
	move   key.Binding
	set    key.Binding
	delete key.Binding
}

type gitLayer struct {
	fetch key.Binding
	push  key.Binding
}

type detailsLayer struct {
	diff    key.Binding
	restore key.Binding
	split   key.Binding
	mark    key.Binding
}

func newKeyMap() keymap {
	bindings := make(map[rune]interface{})
	bindings[' '] = baseLayer{
		abandon:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "abandon")),
		edit:         key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		diffedit:     key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "diff edit")),
		split:        key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "split")),
		squashMode:   key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "squash")),
		rebaseMode:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rebase")),
		bookmarkMode: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "bookmark")),
		gitMode:      key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "git")),
		description:  key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "description")),
		diff:         key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "diff")),
		new:          key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		revset:       key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "revset")),
		refresh:      key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "refresh")),
		undo:         key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "undo")),
		quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	bindings['r'] = rebaseLayer{
		revision: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rebase revision")),
		branch:   key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "rebase branch")),
	}

	bindings['s'] = squashLayer{
		apply: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
	}

	bindings['b'] = bookmarkLayer{
		move:   key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "bookmark move")),
		set:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "bookmark set")),
		delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "bookmark delete")),
	}

	bindings['g'] = gitLayer{
		fetch: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "git fetch")),
		push:  key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "git push")),
	}

	bindings['d'] = detailsLayer{
		diff:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "diff")),
		restore: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restore selected")),
		split:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "split selected")),
		mark:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle selection")),
	}

	return keymap{
		current:  ' ',
		bindings: bindings,
		up:       key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		down:     key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		details:  key.NewBinding(key.WithKeys("l", "details"), key.WithHelp("l", "details")),
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

func (k *keymap) squashMode() {
	k.current = 's'
}

func (k *keymap) bookmarkMode() {
	k.current = 'b'
}

func (k *keymap) resetMode() {
	k.current = ' '
}

func (k *keymap) detailsMode() {
	k.current = 'd'
}

func (k *keymap) ShortHelp() []key.Binding {
	switch b := k.bindings[k.current].(type) {
	case baseLayer:
		return []key.Binding{k.up, k.down, b.revset, b.new, b.edit, b.description, b.diff, b.abandon, b.undo, k.details, b.split, b.squashMode, b.diffedit, b.rebaseMode, b.gitMode, b.bookmarkMode, b.quit}
	case rebaseLayer:
		return []key.Binding{k.up, k.down, b.branch, b.revision, k.cancel}
	case squashLayer:
		return []key.Binding{k.up, k.down, b.apply, k.cancel}
	case gitLayer:
		return []key.Binding{b.push, b.fetch, k.cancel}
	case bookmarkLayer:
		return []key.Binding{b.move, b.set, b.delete, k.cancel}
	case detailsLayer:
		return []key.Binding{k.up, k.down, b.mark, b.diff, b.restore, b.split, k.cancel}
	default:
		if k.current == 'd' {
			return []key.Binding{k.up, k.down, k.cancel}
		}
		return []key.Binding{}
	}
}

func (k *keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
