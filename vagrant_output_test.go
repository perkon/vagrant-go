package vagrant_go

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVagrantOutputLineFromString(t *testing.T) {
	t.Parallel()

	t.Run(
		"with machine readable output that contains an ignored line with `metadata`, it returns nil",
		func(t *testing.T) {
			output := `1546430404,default,metadata,provider,libvirt`

			line := vagrantOutputLineFromString(output)
			assert.Nil(t, line)
		},
	)

	t.Run(
		"with machine readable output that contains an ignored line with `ui`, it returns nil",
		func(t *testing.T) {
			output := `1546430404,default,metadata,provider,libvirt`

			line := vagrantOutputLineFromString(output)
			assert.Nil(t, line)
		},
	)

	t.Run(
		"with machine readable output that contains an ignored line with `action`, it returns nil",
		func(t *testing.T) {
			output := `1546430404,default,metadata,provider,libvirt`

			line := vagrantOutputLineFromString(output)
			assert.Nil(t, line)
		},
	)

	t.Run(
		"with machine readable output that is a parsable line containing single data field, it returns parsed line",
		func(t *testing.T) {
			output := `1546015529,default,box-name,my-debian`

			line := vagrantOutputLineFromString(output)
			require.NotNil(t, line)

			assert.Equal(t, "1546015529", line.timestamp)
			assert.Equal(t, "default", line.target)
			assert.Equal(t, "box-name", line.kind)
			require.Len(t, line.data, 1)
			assert.Equal(t, "my-debian", line.data[0])

		},
	)

	t.Run(
		"with machine readable output that is a parsable line containing multiple data fields, it returns parsed line",
		func(t *testing.T) {
			output := `1546015529,default,box-name,my-debian`

			line := vagrantOutputLineFromString(output)
			require.NotNil(t, line)

			assert.Equal(t, "1546015529", line.timestamp)
			assert.Equal(t, "default", line.target)
			assert.Equal(t, "box-name", line.kind)
			require.Len(t, line.data, 1)
			assert.Equal(t, "my-debian", line.data[0])

		},
	)

	t.Run(
		"with blank output, it returns nil",
		func(t *testing.T) {
			line := vagrantOutputLineFromString("")
			assert.Nil(t, line)
		},
	)
}
