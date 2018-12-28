package vagrant_go

import (
	"testing"
)

func emptyCommandRun(cmd string, args ...string) ([]byte, error) {
	return []byte{}, nil
}

func testClient(t *testing.T, commandRunFunc func(cmd string, args ...string) ([]byte, error)) *Client {
	client, err := NewClient(nil, commandRunFunc)
	if err != nil {
		t.Fatal(err)
	}

	return client
}
