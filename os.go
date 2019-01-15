package vagrant_go

import (
	"os"
)

// Compile-time proof of interface implementation.
var _ OsExecutor = (*osExecutor)(nil)

type OsExecutor interface {
	Chdir(dir string) error
	Getwd() (string, error)
}

type osExecutor struct{}

func (ex *osExecutor) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (ex *osExecutor) Getwd() (string, error) {
	return os.Getwd()
}
