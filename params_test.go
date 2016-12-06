package main

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestParseNoParams(t *testing.T) {
	err, params := parseParams([]string{os.Args[0]})

	fmt.Fprintln(os.Stderr, params)

	assert.Error(t, err)
	assert.NotEqual(t, OK, err)
}


func TestParseHelp(t *testing.T) {
	err, params := parseParams([]string{os.Args[0], "-h"})

	fmt.Fprintln(os.Stderr, params)

	assert.Equal(t, OK, err)
}

func TestParseParams(t *testing.T) {
	err, params := parseParams([]string{os.Args[0], "-t", "pgsql", "-d", "dir", "file1"})

	assert.Nil(t, err)
	assert.NotNil(t, params)

	fmt.Fprintln(os.Stderr, params)

	assert.Equal(t, DUMPTYPE_PGSQL, params.dumptype)

	assert.Equal(t, "file1", params.File())
	assert.Equal(t, "dir", params.destination)
}
