package jj

import (
	"fmt"
	"github.com/idursun/jjui/internal/config"
)

type Commands interface{}

const TEMPLATE = `separate(";", change_id.shortest(1), change_id.shortest(8), coalesce(bookmarks.join(","), "."), current_working_copy, immutable, conflict, empty, hidden, coalesce(author.email(), "(no email set)"), author.timestamp().ago(), commit_id.shortest(1), commit_id.shortest(8), description.first_line())`

type CommandArgs []string

func ConfigGet(key string) CommandArgs {
	return []string{"config", "get", key}
}

func Log(revset string) CommandArgs {
	args := []string{"log", "--color", "never", "--config", "ui.graph.style=curved", "--template", TEMPLATE}
	if revset != "" {
		args = append(args, "-r", revset)
	}
	return args
}

func New(revision string) CommandArgs {
	return []string{"new", "-r", revision}
}

func Edit(changeId string) CommandArgs {
	return []string{"edit", "-r", changeId}
}

func DiffEdit(changeId string) CommandArgs {
	return []string{"diffedit", "-r", changeId}
}

func Split(revision string, files []string) CommandArgs {
	args := []string{"split", "-r", revision}
	args = append(args, files...)
	return args
}

func Describe(revision string) CommandArgs {
	return []string{"describe", "-r", revision, "--edit"}
}

func Abandon(revision string) CommandArgs {
	return []string{"abandon", "-r", revision}
}

func Diff(revision string, fileName string) CommandArgs {
	args := []string{"diff", "-r", revision, "--color", "always"}
	if fileName != "" {
		args = append(args, fileName)
	}
	return args
}

func Restore(revision string, files []string) CommandArgs {
	args := []string{"restore", "-c", revision}
	args = append(args, files...)
	return args
}

func Undo() CommandArgs {
	return []string{"undo"}
}

func Status(revision string) CommandArgs {
	return []string{"log", "-r", revision, "--summary", "--no-graph", "--color", "never", "--template", ""}
}

func BookmarkSet(revision string, name string) CommandArgs {
	return []string{"bookmark", "set", "-r", revision, name}
}

func Squash(from string, destination string) CommandArgs {
	return []string{"squash", "--from", from, "--into", destination}
}

func BookmarkList(revision string) CommandArgs {
	return []string{"bookmark", "list", "-r", fmt.Sprintf("::%s-", revision), "--template", "name ++ if(remote, '@') ++ remote ++ '\n'", "--color", "never"}
}

func BookmarkMove(revision string, bookmark string) CommandArgs {
	return []string{"bookmark", "move", bookmark, "--to", revision}
}

func BookmarkDelete(bookmark string) CommandArgs {
	return []string{"bookmark", "delete", bookmark}
}

func GitFetch(flags ...string) CommandArgs {
	args := []string{"git", "fetch"}
	if flags != nil {
		args = append(args, flags...)
	}
	return args
}

func GitPush(flags ...string) CommandArgs {
	args := []string{"git", "push"}
	if flags != nil {
		args = append(args, flags...)
	}
	return args
}

func Show(revision string) CommandArgs {
	args := []string{"show", "-r", revision, "--color", "always"}
	if config.Current.Preview.ExtraArgs != nil {
		args = append(args, config.Current.Preview.ExtraArgs...)
	}
	return args
}

func Rebase(from string, to string, source string, target string) CommandArgs {
	return []string{"rebase", source, from, target, to}
}

func Evolog(revision string) CommandArgs {
	return []string{"evolog", "-r", revision, "--color", "never", "--template", TEMPLATE, "--config", "ui.graph.style=curved"}
}
