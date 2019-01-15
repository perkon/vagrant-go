package vagrant_go

import (
	"github.com/stretchr/testify/mock"
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

type fakeOsExecutor struct {
	mock.Mock
}

func (f *fakeOsExecutor) Getwd() (string, error) {
	args := f.Called()
	return args.String(0), args.Error(1)
}

func (f *fakeOsExecutor) Chdir(dir string) error {
	args := f.Called(dir)
	return args.Error(0)
}
