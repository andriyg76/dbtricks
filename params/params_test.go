package params

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseNoParams(t *testing.T) {
	err, params := ParseParams([]string{os.Args[0]})

	_, _ = fmt.Fprintln(os.Stderr, params)

	assert.Error(t, err)
	assert.NotEqual(t, OK, err)
}

func TestParseHelp(t *testing.T) {
	err, params := ParseParams([]string{os.Args[0], "-h"})

	_, _ = fmt.Fprintln(os.Stderr, params)

	assert.Equal(t, OK, err)
}

func TestParseParams(t *testing.T) {
	err, params := ParseParams([]string{os.Args[0], "-t", "pgsql", "-d", "dir", "file1"})

	assert.Nil(t, err)
	assert.NotNil(t, params)

	_, _ = fmt.Fprintln(os.Stderr, params)

	assert.Equal(t, Pgsql, params.DumpType)

	assert.Equal(t, "file1", params.File())
	assert.Equal(t, "dir", params.Destination)
}
