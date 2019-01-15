package vagrant_go

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestChdir(t *testing.T) {
	t.Parallel()

	tmpDir, err := ioutil.TempDir("", "example")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	osExecutor := &osExecutor{}
	err = osExecutor.Chdir(tmpDir)
	require.NoError(t, err)

	actualDir, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, actualDir, tmpDir)
}

func TestGetwd(t *testing.T) {
	t.Parallel()

	dir, err := os.Getwd()
	require.NoError(t, err)

	osExecutor := &osExecutor{}
	actualDir, err := osExecutor.Getwd()
	require.NoError(t, err)

	assert.Equal(t, actualDir, dir)
}
