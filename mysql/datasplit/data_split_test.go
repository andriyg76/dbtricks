package datasplit

import (
	"testing"
	"sort"
	"github.com/stretchr/testify/assert"
)

func TestHeadAndValueAscent(t *testing.T) {
	err, value, other := first_value_and_other("11, 'tw\\'oo', ''")
	t.Logf("value=%s tail=%s", value, other)

	assert.NoError(t, err)
	assert.Equal(t, value, "11")
	assert.Equal(t, other, "'tw\\'oo', ''")

	err, value, other = first_value_and_other(other)
	t.Logf("value=%s tail=%s", value, other)
	assert.NoError(t, err)
	assert.Equal(t, value, "tw\\'oo")
	assert.Equal(t, other, "''")

	err, value, other = first_value_and_other(other)
	t.Logf("value=%s tail=%s", value, other)
	assert.NoError(t, err)
	assert.Equal(t, value, "")
	assert.Equal(t, other, "")
}

func TestHeadAndValue(t *testing.T) {
	err, value, other := first_value_and_other("11, 'three'")
	t.Logf("value=%s tail=%s", value, other)

	assert.NoError(t, err)
	assert.Equal(t, value, "11")
	assert.Equal(t, other, "'three'")

	err, value, other = first_value_and_other(other)
	t.Logf("value=%s tail=%s", value, other)
	assert.NoError(t, err)
	assert.Equal(t, value, "three")
	assert.Equal(t, other, "")
}

func TestBufferSort(t *testing.T) {
	one := buffer{
		"11, 'three'",
		"2, 'twoo'",
	}
	sort.Sort(one)

	t.Log("Sorted buffer: ", one)

	assert.Equal(t, buffer{"2, 'twoo'", "11, 'three'"}, one)

	two := buffer{
		"11, 'three'",
		"11, 'tw\\'oo'",
	}

	sort.Sort(two)

	t.Log("Sorted buffer: ", two)
	assert.Equal(t, buffer{"11, 'three'", "11, 'tw\\'oo'"}, two)
}