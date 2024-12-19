[![Build & Test](https://github.com/idursun/jjui/actions/workflows/go.yml/badge.svg)](https://github.com/idursun/jjui/actions/workflows/go.yml)

# Jujutsu UI

A TUI for working with [Jujutsu version control system](https://github.com/martinvonz/jj).

I have built `jjui` according to my own needs and will keep adding new features as I need them.

## Features

Currently, you can:

### Change revset with auto-complete
You can change revset while enjoying auto-complete and signature help while typing.

![GIF](./docs/jjui_revset.gif)

### Description
You can edit or update description of a revision.

![GIF](./docs/jjui_description.gif)

### Abandon
You can abandon a revision.

![GIF](./docs/jjui_abandon.gif)

### Rebase
You can rebase a revision or a branch onto another revision in the revision tree.

![GIF](./docs/jjui_rebase.gif)
### Bookmarks
You can move bookmarks to the revision you selected.

![GIF](./docs/jjui_bookmarks.gif)

### Diffs
You can see diffs of revisions.

![GIF](./docs/jjui_diff.gif)

### Split
You can split revisions.

![GIF](./docs/jjui_split.gif)
Additionally,

* Create a _new_ revision
* _Edit_ a revision
* Git _push_/_fetch_

## Installation

## Nix

You can install `jjui` using nix from the unstable channel.

```shell
nix-env -iA nixpkgs.jjui
```

## From source

You can build `jjui` from source.

```
git clone https://github.com/idursun/jjui.git
cd jjui
go install cmd/jjui.go
```



## Compatibility

It's compatible with jj **v0.21**+.

### Contributing

Feel free to submit a pull request.
