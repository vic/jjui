package revset

import (
	"fmt"
	"strings"
	"unicode"
)

type FunctionDefinition struct {
	Name          string
	HasParameters bool
	SignatureHelp string
	IsAlias       bool
}

var AllFunctions = []FunctionDefinition{
	{"all", false, "all(): All commits", false},
	{"mine", false, "mine(): Your own commits", false},
	{"empty", false, "empty(): The empty set", false},
	{"trunk", false, "trunk(): The trunk of the repository", false},
	{"root", false, "root(): The root commit", false},
	{"description", true, "description(pattern): Commits that have a description matching the given string pattern", false},
	{"author", true, "author(pattern): Commits with the author's name or email matching the given string pattern", false},
	{"author_date", true, "author_date(pattern): Commits with author dates matching the specified date pattern.", false},
	{"committer", true, "committer(pattern): Commits with the committer's name or email matching the given pattern", false},
	{"committer_date", true, "committer_date(pattern): Commits with committer dates matching the specified date pattern", false},
	{"tags", true, "tags([pattern]): All tag targets. If pattern is specified, this selects the tags whose name match the given string pattern", false},
	{"files", true, "files(expression): Commits modifying paths matching the given fileset expression", false},
	{"latest", true, "latest(x[, count]): Latest count commits in x", false},
	{"bookmarks", true, "bookmarks([pattern]): If pattern is specified, this selects the bookmarks whose name match the given string pattern", false},
	{"conflicts", false, "conflicts(): Commits with conflicts", false},
	{"diff_contains", true, "diff_contains(text[, files]): Commits containing the given text in their diffs", false},
	{"descendants", true, "descendants(x[, depth]): Returns the descendants of x limited to the given depth", false},
	{"parents", true, "parents(x): Same as x-", false},
	{"ancestors", true, "ancestors(x[, depth]): Returns the ancestors of x limited to the given depth", false},
	{"connected", true, "connected(x): Same as x::x. Useful when x includes several commits", false},
	{"git_head", false, "git_head(): The commit referred to by Git's HEAD", false},
	{"git_refs", false, "git_refs(): All Git refs", false},
	{"heads", true, "heads(x): Commits in x that are not ancestors of other commits in x", false},
	{"fork_point", true, "fork_point(x): The fork point of all commits in x", false},
	{"merges", true, "merges(x): Commits in x with more than one parent", false},
	{"remote_bookmarks", true, "remote_bookmarks([bookmark_pattern[, [remote=]remote_pattern]]): All remote bookmarks targets across all remotes", false},
	{"present", true, "present(x): Same as x, but evaluated to none() if any of the commits in x doesn't exist", false},
	{"coalesce", true, "coalesce(revsets...): Commits in the first revset in the list of revsets which does not evaluate to none()", false},
	{"working_copies", false, "working_copies(): All working copies", false},
	{"at_operation", true, "at_operation(op, x): Evaluates to x at the specified operation", false},
	{"tracked_remote_bookmarks", true, "tracked_remote_bookmarks([bookmark_pattern[, [remote=]remote_pattern]])", false},
	{"untracked_remote_bookmarks", true, "untracked_remote_bookmarks([bookmark_pattern[, [remote=]remote_pattern]])", false},
	{"visible_heads", false, "visible_heads(): All visible heads in the repo", false},
	{"reachable", true, "reachable(srcs, domain): All commits reachable from srcs within domain, traversing all parent and child edges", false},
	{"roots", true, "roots(x): Commits in x that are not descendants of other commits in x", false},
	{"children", true, "children(x): Same as x+", false},
}

func GetFunctionByName(name string) *FunctionDefinition {
	for _, fn := range AllFunctions {
		if fn.Name == name {
			return &fn
		}
	}
	return nil
}

type CompletionProvider struct {
}

func NewCompletionProvider(aliases map[string]string) *CompletionProvider {
	for alias, expansion := range aliases {
		hasParameters := false
		signatureHelp := fmt.Sprintf("%s: %s", alias, expansion)

		if strings.Index(alias, "(") < strings.LastIndex(alias, ")") {
			hasParameters = true
			alias = alias[:strings.Index(alias, "(")]
		} else if strings.HasSuffix(alias, "()") {
			hasParameters = false
			alias = alias[:len(alias)-2]
		}

		AllFunctions = append(AllFunctions, FunctionDefinition{
			Name:          alias,
			HasParameters: hasParameters,
			SignatureHelp: signatureHelp,
			IsAlias:       true,
		})
	}
	return &CompletionProvider{}
}

func (p *CompletionProvider) GetCompletions(input string) []string {
	var suggestions []string
	if input == "" {
		for _, function := range AllFunctions {
			if !function.IsAlias {
				continue
			}
			suggestions = append(suggestions, function.Name)
		}
		return suggestions
	}

	lastToken := getLastToken(input)
	if lastToken == "" {
		return nil
	}

	for _, fn := range AllFunctions {
		if strings.HasPrefix(fn.Name, lastToken) {
			suggestions = append(suggestions, fn.Name)
		}
	}

	return suggestions
}

func (p *CompletionProvider) GetSignatureHelp(input string) string {
	helpFunction := extractLastFunctionName(input)
	if helpFunction == "" {
		return ""
	}

	if fn := GetFunctionByName(helpFunction); fn != nil {
		return fn.SignatureHelp
	}

	return ""
}

func (p *CompletionProvider) GetLastToken(input string) (int, string) {
	lastIndex := strings.LastIndexFunc(input, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == '|' || r == '&' || r == '~' || r == '(' || r == '.' || r == ':'
	})

	if lastIndex == -1 {
		return 0, input
	}

	return lastIndex + 1, input[lastIndex+1:]
}

func extractLastFunctionName(input string) string {
	lastOpenParen := strings.LastIndex(input, "(")
	if lastOpenParen == -1 {
		return ""
	}

	parenCount := 1
	for i := lastOpenParen + 1; i < len(input); i++ {
		if input[i] == '(' {
			parenCount++
		} else if input[i] == ')' {
			parenCount--
		}

		if parenCount == 0 && i+1 < len(input) {
			for j := i + 1; j < len(input); j++ {
				ch := input[j]
				if ch == '|' || ch == '&' || ch == ',' || !unicode.IsSpace(rune(ch)) {
					return ""
				}

				if !unicode.IsSpace(rune(ch)) {
					break
				}
			}
			break
		}
	}

	startIndex := lastOpenParen
	for startIndex > 0 {
		startIndex--
		if !isValidFunctionNameChar(rune(input[startIndex])) {
			startIndex++ // Move back to the valid character
			break
		}
	}

	if startIndex <= lastOpenParen {
		funcName := input[startIndex:lastOpenParen]
		return funcName
	}

	return ""
}

func isValidFunctionNameChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func getLastToken(input string) string {
	lastIndex := strings.LastIndexFunc(input, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == '|' || r == '&' || r == '~' || r == '(' || r == '.' || r == ':'
	})

	if lastIndex == -1 {
		return input
	}

	if lastIndex+1 < len(input) {
		return input[lastIndex+1:]
	}

	return ""
}
