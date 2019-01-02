package vagrant_go

import (
	"testing"
)

func emptyCommandRunFunc(cmd string, args ...string) ([]byte, error) {
	return []byte{}, nil
}

func emptyLookPathFunc(file string) (string, error) {
	return "", nil
}

func testClient(
	t *testing.T,
	commandRunFunc func(cmd string, args ...string) ([]byte, error),
	lookPathFunc func(file string) (string, error),
) *Client {
	client, err := NewClient(nil, commandRunFunc, lookPathFunc)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func emptyTestClient(t *testing.T) *Client {
	client, err := NewClient(nil, emptyCommandRunFunc, emptyLookPathFunc)
	if err != nil {
		t.Fatal(err)
	}

	return client
}
