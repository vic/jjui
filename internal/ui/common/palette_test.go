package common

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"github.com/stretchr/testify/assert"
)

const (
	Black  = "0"
	Red    = "1"
	Green  = "2"
	Yellow = "3"
	Blue   = "4"
	Cyan   = "6"
	White  = "7"
)

func TestPalette_Get(t *testing.T) {
	// Set up a palette with test styles using the add method
	p := NewPalette()

	// Add styles using the palette's add method
	p.add("text", lipgloss.NewStyle().Foreground(lipgloss.Color(White)))
	p.add("selected", lipgloss.NewStyle().Background(lipgloss.Color(Black)).Bold(true))
	p.add("revisions", lipgloss.NewStyle().Italic(true))
	p.add("revisions text", lipgloss.NewStyle().Foreground(lipgloss.Color(Cyan)).Background(lipgloss.Color(Green)))

	tests := []struct {
		name     string
		selector string
		want     lipgloss.Style
		palette  *Palette
	}{
		{
			name:     "exact match for single label",
			selector: "text",
			want:     lipgloss.NewStyle().Foreground(lipgloss.Color(White)),
			palette:  p,
		},
		{
			name:     "combined labels",
			selector: "revisions selected",
			want:     lipgloss.NewStyle().Background(lipgloss.Color(Black)).Bold(true).Italic(true),
			palette:  p,
		},
		{
			name:     "non-existent label",
			selector: "nonexistent",
			want:     lipgloss.NewStyle(),
			palette:  p,
		},
		{
			name:     "mixed existing and non-existent labels",
			selector: "text nonexistent",
			want:     lipgloss.NewStyle().Foreground(lipgloss.Color(White)),
			palette:  p,
		},
		{
			name:     "empty selector",
			selector: "",
			want:     lipgloss.NewStyle(),
			palette:  p,
		},
		{
			name:     "exact match for compound label",
			selector: "revisions text",
			want:     lipgloss.NewStyle().Foreground(lipgloss.Color(Cyan)).Background(lipgloss.Color(Green)).Italic(true),
			palette:  p,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.palette.Get(tt.selector)

			// Compare foreground colours
			assert.Equal(t, tt.want.GetForeground(), got.GetForeground(), "foreground color mismatch")

			// Compare background colours
			assert.Equal(t, tt.want.GetBackground(), got.GetBackground(), "background color mismatch")

			// Compare style attributes
			assert.Equal(t, tt.want.GetBold(), got.GetBold(), "bold attribute mismatch")
			assert.Equal(t, tt.want.GetItalic(), got.GetItalic(), "italic attribute mismatch")
		})
	}
}

func TestPalette_Update(t *testing.T) {
	tests := []struct {
		name     string
		styleMap map[string]config.Color
		selector string
		want     lipgloss.Style
	}{
		{
			name: "basic color update",
			styleMap: map[string]config.Color{
				"text": {Fg: Red},
			},
			selector: "text",
			want:     lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		},
		{
			name: "update with multiple attributes",
			styleMap: map[string]config.Color{
				"heading": {Fg: Blue, Bold: true, Italic: true},
			},
			selector: "heading",
			want:     lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true).Italic(true),
		},
		{
			name: "update with background color",
			styleMap: map[string]config.Color{
				"highlight": {Fg: Black, Bg: Yellow},
			},
			selector: "highlight",
			want:     lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("3")),
		},
		{
			name: "diff shortcuts",
			styleMap: map[string]config.Color{
				"diff added":    {Fg: Green},
				"diff renamed":  {Fg: Blue},
				"diff modified": {Fg: Yellow},
				"diff removed":  {Fg: Red},
			},
			selector: "added",
			want:     lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPalette()
			p.Update(tt.styleMap)

			got := p.Get(tt.selector)

			if tt.name == "diff shortcuts" {
				// Check that all diff shortcuts were properly added
				assert.Equal(t, lipgloss.Color("2"), p.Get("added").GetForeground(), "added style not set correctly")
				assert.Equal(t, lipgloss.Color("4"), p.Get("renamed").GetForeground(), "renamed style not set correctly")
				assert.Equal(t, lipgloss.Color("3"), p.Get("modified").GetForeground(), "modified style not set correctly")
				assert.Equal(t, lipgloss.Color("1"), p.Get("deleted").GetForeground(), "deleted style not set correctly")
			} else {
				assert.Equal(t, tt.want.GetForeground(), got.GetForeground(), "foreground color mismatch")
				assert.Equal(t, tt.want.GetBackground(), got.GetBackground(), "background color mismatch")
				assert.Equal(t, tt.want.GetBold(), got.GetBold(), "bold attribute mismatch")
				assert.Equal(t, tt.want.GetItalic(), got.GetItalic(), "italic attribute mismatch")
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  lipgloss.Color
	}{
		{
			name:  "hex color",
			color: "#ff0000",
			want:  lipgloss.Color("#ff0000"),
		},
		{
			name:  "ansi256 color by number",
			color: "123",
			want:  lipgloss.Color("123"),
		},
		{
			name:  "named color - red",
			color: "red",
			want:  lipgloss.Color("1"),
		},
		{
			name:  "named color - bright blue",
			color: "bright blue",
			want:  lipgloss.Color("12"),
		},
		{
			name:  "ansi-color prefix",
			color: "ansi-color-42",
			want:  lipgloss.Color("42"),
		},
		{
			name:  "invalid color",
			color: "not-a-color",
			want:  lipgloss.Color(""),
		},
		{
			name:  "out of range ansi256",
			color: "300",
			want:  lipgloss.Color(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseColor(tt.color)
			assert.Equal(t, tt.want, got)
		})
	}
}
