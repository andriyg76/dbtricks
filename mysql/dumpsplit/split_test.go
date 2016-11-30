package dumpsplit
import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestRegexps(t *testing.T) {
	match, table_name := match_table_structure("-- Table structure for table `django_site`")

	assert.True(t, match)
	assert.Equal(t, "django_site", table_name)

	match, insert_into, table_data := match_table_data(
		"INSERT INTO `django_site` (`id`, `domain`, `name`) VALUES " +
			"(1,'http://caritas-kolomyya.org','caritas');")

	assert.True(t, match)
	assert.Equal(t, insert_into, "INSERT INTO `django_site` (`id`, `domain`, `name`) VALUES")
	assert.Equal(t, table_data, "1,'http://caritas-kolomyya.org','caritas'")
}

