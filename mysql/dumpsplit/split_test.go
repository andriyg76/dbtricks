package dumpsplit
import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestRegexps(t *testing.T) {
	match, tableName := matchTableStructure("-- Table structure for table `django_site`")

	assert.True(t, match)
	assert.Equal(t, "django_site", tableName)

	match, insertInto, tableData := matchTableData(
		"INSERT INTO `django_site` (`id`, `domain`, `name`) VALUES " +
			"(1,'http://caritas-kolomyya.org','caritas');")

	assert.True(t, match)
	assert.Equal(t, insertInto, "INSERT INTO `django_site` (`id`, `domain`, `name`) VALUES")
	assert.Equal(t, tableData, "1,'http://caritas-kolomyya.org','caritas'")
}

