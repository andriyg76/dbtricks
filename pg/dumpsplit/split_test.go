package dumpsplit
import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestRegexps(t *testing.T) {
	match := matchToCopy("COPY table (column1, column2, column2) FROM stdin;")

	assert.True(t, match)

	match, table, schema := matchToDataComment(
		"-- Data for Name: table; Type: TABLE DATA; Schema: s1;")

	assert.True(t, match)
	assert.Equal(t, table, "table")
	assert.Equal(t, schema, "s1")

	match = matchToConstraintComment("-- Name: PK_table; Type: CONSTRAINT; Schema: s1")
	assert.True(t, match)

	match = matchToConstraintComment("-- Name: IX_TABLE_COLUMN1; Type: INDEX; Schema: s1")
	assert.True(t, match)

	match = matchToConstraintComment("-- Name: IX_TABLE_COLUMN1; Type: SEQUENCE; Schema: s1")
	assert.False(t, match)
}

