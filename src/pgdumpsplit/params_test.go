package pgdumpsplit

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestParseParams(t *testing.T) {
	params, err := ParseParams([]string{os.Args[0], "-h"})

	fmt.Fprintln(os.Stderr, params.String())
	assert.Equal(t, err, nil)

	//assert.Equal(t, true, params.IsHelp())
}
