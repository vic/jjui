package jj

import (
	"strings"
)

const moveBookmarkTemplate = `separate(";", if(remote, name ++ "@" ++ remote, name), if(remote, "true", "false"), tracked, conflict, normal_target.contained_in("%s"), normal_target.commit_id().shortest(1)) ++ "\n"`
const allBookmarkTemplate = `separate(";", if(remote, name ++ "@" ++ remote, name), if(remote, "true", "false"), tracked, conflict, 'false', normal_target.commit_id().shortest(1)) ++ "\n"`

type Bookmark struct {
	Name      string
	Tracked   bool
	Remote    bool
	Conflict  bool
	Backwards bool
	CommitId  string
}

func ParseBookmarkListOutput(output string) []Bookmark {
	bookmarks := strings.Split(output, "\n")
	var result []Bookmark
	for _, b := range bookmarks {
		parts := strings.Split(b, ";")
		if len(parts) < 5 {
			continue
		} else {
			bookmark := Bookmark{
				Name:      parts[0],
				Remote:    parts[1] == "true",
				Tracked:   parts[2] == "true",
				Conflict:  parts[3] == "true",
				Backwards: parts[4] == "true",
				CommitId:  parts[5],
			}
			result = append(result, bookmark)
		}
	}
	return result

}
