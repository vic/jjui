package test

import (
    "github.com/stretchr/testify/assert"
    "os"
    "testing"

    "jjui/internal/jj"
)

func Test_Parse_Tree(t *testing.T) {
	fileName := "testdata/many-levels.log"
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	defer file.Close()

	node := jj.Parse(file)
	treeRenderer := jj.NewTreeRenderer()
	buffer := treeRenderer.RenderTree(node)
	content, err := os.ReadFile("testdata/many-levels.rendered")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	assert.Equal(t, string(content), buffer)
}