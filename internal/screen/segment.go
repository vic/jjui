package screen

import (
	"fmt"
	"strconv"
	"strings"
)

type Segment struct {
	Text   string
	Params string
}

func (s Segment) String() string {
	if s.Text == "\n" {
		return s.Text
	}
	if s.Params == "" {
		return s.Text
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", s.Params, s.Text)
}

func (s Segment) WithBackground(bg string) string {
	var newParts []string
	parts := strings.Split(s.Params, ";")

	i := 0
	for i < len(parts) {
		part := parts[i]
		num, err := strconv.Atoi(part)
		if err != nil {
			i++
			continue
		}
		p := num

		isBg := false
		if (p >= 40 && p <= 49) || (p >= 100 && p <= 109) {
			isBg = true
		}

		if !isBg {
			newParts = append(newParts, part)
		}
		i++
	}

	for _, part := range strings.Split(bg, ";") {
		if _, err := strconv.Atoi(part); err != nil {
			panic(fmt.Sprintf("invalid background parameter %q", part))
		}
		newParts = append(newParts, part)
	}

	return Segment{
		Text:   s.Text,
		Params: strings.Join(newParts, ";"),
	}.String()
}

func (s Segment) StyleEqual(other Segment) bool {
	return s.Params == other.Params
}
