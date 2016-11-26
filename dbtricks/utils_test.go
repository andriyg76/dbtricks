package dbtricks
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConvertInt(t *testing.T) {
	val := IntInBase(0, 36, 0)

	assert.Equal(t, "0", val)

	val = IntInBase(36, 36, 0)

	assert.Equal(t, "10", val)

	val = IntInBase(36, 36, 4)

	assert.Equal(t, "0010", val)

	val = IntInBase(-36, 36, 4)

	assert.Equal(t, "-0010", val)
}

