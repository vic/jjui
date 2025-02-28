[![Build & Test](https://github.com/idursun/jjui/actions/workflows/go.yml/badge.svg)](https://github.com/idursun/jjui/actions/workflows/go.yml)

# Jujutsu UI

`jjui` is a terminal user interface for working with [Jujutsu version control system](https://github.com/jj-vcs/jj). I have built it according to my own needs and will keep adding new features as I need them. I am open to feature requests and contributions.

## Features

Currently, you can:

### Change revset with auto-complete
You can change revset while enjoying auto-complete and signature help while typing.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_revset.gif)

### Rebase
You can rebase a revision or a branch onto another revision in the revision tree.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_rebase.gif)

### Squash
You can squash revisions into one revision, by pressing `S`. The following revision will be automatically selected. However, you can change the selection by using `j` and `k`.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_squash.gif)

### Show revision details

Pressing `l` (as in going right into details of a revision) will open the details view of the revision you selected.

You can use `d` to show the diff of the each highlighted file.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_details_diff.gif)

#### Splitting files in a revision
You can use `s` to split the revision based on the files you have selected. You can select by pressing `space`. If no file is selected, currently highlighted file will be split.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_details_split.gif)

### Restoring files in a revision
You can use `r` to restore the selected file (i.e. discard the changes).

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_details_restore.gif)

### Description
You can edit or update description of a revision.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_description.gif)

### Abandon
You can abandon a revision.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_abandon.gif)

### Bookmarks
You can move bookmarks to the revision you selected.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_bookmarks.gif)

### Split
You can split revisions.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_split.gif)

### Preview
You can open the preview window by pressing `p`. Preview window displays output of the `jj show` command of the selected revision. If the currenlty selected item is a file, then `jj diff` output is displayed. 

While the preview window is showing, you can press; `ctrl+n` to scroll one line down, `ctrl+p` to scroll one line up, `ctrl+n` to scroll half page down, `ctrl+u` to scroll half page up. 

Additionally, you can press `tab` to focus in and out of the preview window. Once in the focus mode, you can use normal (e.g. `j`, `k`, `d`, `u`) navigation keys as they are not bound to the revision tree view. 

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_preview.gif)

### Diffs
You can see diffs of revisions.

![GIF](https://github.com/idursun/jjui/wiki/gifs/jjui_diff.gif)

Additionally,

* Create a _new_ revision
* _Edit_ a revision
* Git _push_/_fetch_
* Undo last change

## Installation

### Nix

You can install `jjui` using nix from the unstable channel.

```shell
nix-env -iA nixpkgs.jjui
```

### From source

You can build `jjui` from source.

```
git clone https://github.com/idursun/jjui.git
cd jjui
go install cmd/jjui/main.go
```

### From go install
```
go install github.com/idursun/jjui/cmd/jjui@latest
```

### From pre-built binaries
You can download pre-built binaries from the [releases](https://github.com/idursun/jjui/releases) page.

## Compatibility

It's compatible with jj **v0.21**+.

## Contributing

Feel free to submit a pull request.
