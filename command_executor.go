package vagrant_go

import (
	"os/exec"
)

func realCommandRunFunc(cmd string, args ...string) ([]byte, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	return out, err
}

func realLookPathFunc(file string) (string, error) {
	path, err := exec.LookPath(file)
	return path, err
}
