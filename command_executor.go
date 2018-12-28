package vagrant_go

import (
	"os/exec"
)

func realCommandFunc(cmd string, args ...string) ([]byte, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	return out, err
}
