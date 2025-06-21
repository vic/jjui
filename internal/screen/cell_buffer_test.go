package screen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStacked_OverlappingDoubleWidth(t *testing.T) {
	first := "🤬."
	stacked := Stacked(first, "|", 1, 0)
	assert.Equal(t, " |.", stacked)
}
