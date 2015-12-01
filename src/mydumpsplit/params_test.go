package mydumpsplit
import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseParams(t *testing.T) {
	_, err := ParseParams([]string{os.Args[0], "-h"})

	assert.Equal(t, err, nil)

	//assert.Equal(t, true, params.IsHelp())
}
