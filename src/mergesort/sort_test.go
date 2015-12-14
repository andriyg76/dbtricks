package mergesort
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type Strin Strings  {

}
func TestNofileOpen(t *testing.T) {
	err, _ := NewAsyncFileReader(nil)

	t.Log(err)

	assert.NotNil(t, err)

	err, reader := NewAsyncFileReader(arrayReader)
}