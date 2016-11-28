package main

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestParseParams(t *testing.T) {
	params:= parseParams([]string{os.Args[0], "-h", "file1"})

	fmt.Fprintln(os.Stderr, params)
	assert.Equal(t, params.Error(), nil)

	assert.Equal(t, true, params.IsHelp())

	assert.Equal(t, "file1", params.File())
}
