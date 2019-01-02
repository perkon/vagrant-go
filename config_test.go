package vagrant_go

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultBinaryName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "vagrant", defaultBinaryName)
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	config := DefaultConfig()
	require.NotNil(t, config)
	assert.Equal(t, defaultBinaryName, config.BinaryName)
}
