package vagrant_go

import (
	"github.com/palantir/stacktrace"
)

// Compile-time proof of interface implementation.
var _ BoxAPI = (*boxAPI)(nil)

type BoxAPI interface {
	List() ([]*Box, error)
}

type boxAPI struct {
	client *Client
}

type Box struct {
	Name     string
	Provider string
	Version  string
}

func (api *boxAPI) List() ([]*Box, error) {
	outputLines, err := api.client.executeVagrantCommand("box", "list")
	if err != nil {
		return nil, stacktrace.Propagate(err, "command execution failed")
	}

	var name, provider, version string

	// NOTE: Use 0 element slice in case there's nothing to return
	//noinspection GoPreferNilSlice
	boxes := []*Box{}

	for _, line := range outputLines {
		switch line.kind {
		case "box-name":
			if len(name) > 0 {
				boxes = append(
					boxes,
					&Box{
						Name:     name,
						Provider: provider,
						Version:  version,
					},
				)
			}

			name = line.data[0]
			provider = ""
			version = ""
		case "box-provider":
			provider = line.data[0]
		case "box-version":
			version = line.data[0]
		}
	}

	// NOTE: Add last box from output lines
	if len(name) > 0 {
		boxes = append(
			boxes,
			&Box{
				Name:     name,
				Provider: provider,
				Version:  version,
			},
		)
	}

	return boxes, nil
}
