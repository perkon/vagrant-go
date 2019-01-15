package vagrant_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	t.Run(
		"with array containing 'needle', it returns true",
		func(t *testing.T) {
			t.Parallel()

			array := []string{"foo", "bar"}

			actual := contains(array, "foo")
			assert.True(t, actual)
		},
	)

	t.Run(
		"with array not containing 'needle', it returns false",
		func(t *testing.T) {
			t.Parallel()

			array := []string{"foo", "bar"}

			actual := contains(array, "foobar")
			assert.False(t, actual)
		},
	)
}
