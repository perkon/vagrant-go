package vagrant_go

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoxApiList(t *testing.T) {
	t.Parallel()

	t.Run(
		"with no vagrant boxes available",
		func(t *testing.T) {
			boxAPI := &boxAPI{
				client: testClient(
					t,
					emptyCommandRun,
				),
			}

			boxes, err := boxAPI.List()

			require.NoError(t, err)
			assert.Empty(t, boxes)
		},
	)

	t.Run(
		"with 1 vagrant box available",
		func(t *testing.T) {
			singleBoxCommandRunFunc := func(cmd string, args ...string) ([]byte, error) {
				output := `
1546015529,,ui,info,my-debian (libvirt%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-debian
1546015529,,box-provider,libvirt
1546015529,,box-version,1.2.3
`
				return []byte(output), nil
			}

			boxAPI := &boxAPI{
				client: testClient(
					t,
					singleBoxCommandRunFunc,
				),
			}

			boxes, err := boxAPI.List()

			require.NoError(t, err)
			assert.Len(t, boxes, 1)

			assert.Equal(t, boxes[0].Name, "my-debian")
			assert.Equal(t, boxes[0].Provider, "libvirt")
			assert.Equal(t, boxes[0].Version, "1.2.3")
		},
	)

	t.Run(
		"with 3 vagrant boxes available",
		func(t *testing.T) {
			singleBoxCommandRunFunc := func(cmd string, args ...string) ([]byte, error) {
				output := `
1546015529,,ui,info,my-debian (libvirt%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-debian
1546015529,,box-provider,libvirt
1546015529,,box-version,1.2.3
1546015529,,ui,info,my-vbox-debian (virtualbox%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-vbox-debian
1546015529,,box-provider,virtualbox
1546015529,,box-version,1.2.4
1546015529,,ui,info,my-vmware-debian (vmware%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-vmware-debian
1546015529,,box-provider,vmware
1546015529,,box-version,1.2.5
`
				return []byte(output), nil
			}

			boxAPI := &boxAPI{
				client: testClient(
					t,
					singleBoxCommandRunFunc,
				),
			}

			boxes, err := boxAPI.List()

			require.NoError(t, err)
			assert.Len(t, boxes, 3)

			assert.Equal(t, boxes[0].Name, "my-debian")
			assert.Equal(t, boxes[0].Provider, "libvirt")
			assert.Equal(t, boxes[0].Version, "1.2.3")

			assert.Equal(t, boxes[1].Name, "my-vbox-debian")
			assert.Equal(t, boxes[1].Provider, "virtualbox")
			assert.Equal(t, boxes[1].Version, "1.2.4")

			assert.Equal(t, boxes[2].Name, "my-vmware-debian")
			assert.Equal(t, boxes[2].Provider, "vmware")
			assert.Equal(t, boxes[2].Version, "1.2.5")
		},
	)
}
