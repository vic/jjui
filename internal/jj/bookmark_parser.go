package jj

import (
	"strings"
)

const moveBookmarkTemplate = `separate(";", name, if(remote, "remote", "."), tracked, conflict, normal_target.contained_in("%s"), normal_target.commit_id().shortest(1)) ++ "\n"`
const allBookmarkTemplate = `separate(";", name, if(remote, remote, "."), tracked, conflict, 'false', normal_target.commit_id().shortest(1)) ++ "\n"`

type BookmarkRemote struct {
	Remote   string
	CommitId string
	Tracked  bool
}

type Bookmark struct {
	Name      string
	Remotes   []BookmarkRemote
	Conflict  bool
	Backwards bool
	CommitId  string
}

func (b Bookmark) IsLocal() bool {
	return len(b.Remotes) == 0
}

func ParseBookmarkListOutput(output string) []Bookmark {
	lines := strings.Split(output, "\n")
	var bookmarks []Bookmark
	for _, b := range lines {
		parts := strings.Split(b, ";")
		if len(parts) < 5 {
			continue
		} else {
			name := parts[0]
			remoteName := parts[1]
			tracked := parts[2] == "true"
			conflict := parts[3] == "true"
			backwards := parts[4] == "true"
			commitId := parts[5]
			if remoteName == "." {
				bookmark := Bookmark{
					Name:      name,
					Conflict:  conflict,
					Backwards: backwards,
					CommitId:  commitId,
				}
				bookmarks = append(bookmarks, bookmark)
			} else if len(bookmarks) > 0 {
				previous := &bookmarks[len(bookmarks)-1]
				remote := BookmarkRemote{
					Remote:   remoteName,
					Tracked:  tracked,
					CommitId: commitId,
				}
				if remoteName == "origin" && len(previous.Remotes) > 0 {
					// add the origin remote to the front of the list
					previous.Remotes = append([]BookmarkRemote{remote}, previous.Remotes...)
				} else {
					previous.Remotes = append(previous.Remotes, remote)
				}
			}
		}
	}
	return bookmarks
}
