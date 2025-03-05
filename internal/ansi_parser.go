package ui

import (
	"fmt"
	"strconv"
	"strings"
)

type Segment struct {
	Text   string
	Params []int
}

func (s Segment) String() string {
	if len(s.Params) == 0 {
		return s.Text
	}

	params := make([]string, len(s.Params))
	for i, p := range s.Params {
		params[i] = strconv.Itoa(p)
	}
	return fmt.Sprintf(
		"\x1b[%sm%s\x1b[0m",
		strings.Join(params, ";"),
		s.Text,
	)
}

func (s Segment) StyleEqual(other Segment) bool {
	if len(other.Params) != len(s.Params) {
		return false
	}
	for i, p := range s.Params {
		if p != other.Params[i] {
			return false
		}
	}
	return true
}

func Parse(raw string) []Segment {
	var segments []Segment
	var buffer strings.Builder
	var params []int
	pos := 0

	for pos < len(raw) {
		if raw[pos] == 0x1B && pos+1 < len(raw) && raw[pos+1] == '[' {
			// Save current buffer
			if buffer.Len() > 0 {
				segments = append(segments, Segment{
					Text:   buffer.String(),
					Params: params,
				})
				params = nil
				buffer.Reset()
			}

			// Extract full escape sequence
			end := strings.IndexByte(raw[pos:], 'm')
			if end == -1 {
				pos++
				continue
			}
			end += pos

			// Parse parameters
			seq := raw[pos+2 : end]
			for _, param := range strings.Split(seq, ";") {
				if param == "" {
					continue
				}
				if num, err := strconv.Atoi(param); err == nil {
					params = append(params, num)
				}
			}

			// Handle reset
			if len(params) == 1 && params[0] == 0 {
				params = nil
			}

			pos = end + 1
		} else {
			buffer.WriteByte(raw[pos])
			pos++
		}
	}

	// Add remaining text
	if buffer.Len() > 0 {
		segments = append(segments, Segment{
			Text:   buffer.String(),
			Params: params,
		})
	}

	return segments
}
