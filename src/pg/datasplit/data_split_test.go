package datasplit

import (
	"testing"
	"sort"
	"github.com/stretchr/testify/assert"
)


func TestBufferSort(t *testing.T) {
	one := buffer{
		"11\tthree",
		"2\ttwoo",
	}
	sort.Sort(one)

	t.Log("Sorted buffer: ", one)

	assert.Equal(t, buffer{"2\ttwoo", "11\tthree"}, one)

	two := buffer{
		"11\tthree",
		"11\ttwoo",
	}

	sort.Sort(two)

	t.Log("Sorted buffer: ", two)
	assert.Equal(t, buffer{"11\tthree", "11\ttwoo"}, two)
}