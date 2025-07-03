package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	content := `
[ui]
highlight_light = "#a0a0a0"
`
	config := &Config{}
	err := config.Load(content)
	assert.NoError(t, err)
	assert.Equal(t, "#a0a0a0", config.UI.HighlightLight)
}
