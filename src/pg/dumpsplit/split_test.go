package dumpsplit
import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestRegexps(t *testing.T) {
	match := match_to_copy("COPY table (column1, column2, column2) FROM stdin;")

	assert.True(t, match)

	match, table, schema := match_to_data_comment(
		"-- Data for Name: table; Type: TABLE DATA; Schema: s1;")

	assert.True(t, match)
	assert.Equal(t, table, "table")
	assert.Equal(t, schema, "s1")

	match = match_to_constraint_comment("-- Name: PK_table; Type: CONSTRAINT; Schema: s1")
	assert.True(t, match)

	match = match_to_constraint_comment("-- Name: IX_TABLE_COLUMN1; Type: INDEX; Schema: s1")
	assert.True(t, match)

	match = match_to_constraint_comment("-- Name: IX_TABLE_COLUMN1; Type: SEQUENCE; Schema: s1")
	assert.False(t, match)
}

