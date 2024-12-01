package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"jjui/internal/jj"
)

func Test_parseLogOutput_ManyLevels(t *testing.T) {
	fileName := "testdata/many-levels.log"
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	rows := jj.Parse(file)
	assert.Equal(t, 8, len(rows))
}

func Test_parseLogOutput_TwoLevels(t *testing.T) {
	fileName := "testdata/two-level.log"
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	rows := jj.Parse(file)
	assert.Equal(t, 10, len(rows))
}
func Test_Parse_ElidedRevisions(t *testing.T) {
	fileName := "testdata/elided-revisions.log"
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	rows := jj.Parse(file)
	assert.Equal(t, 6, len(rows))
}
