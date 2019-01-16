package vagrant_go

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func realCommandRunFunc(cmd string, args ...string) ([]byte, error) {
	var outBuffer bytes.Buffer

	execCmd := exec.Command(cmd, args...)

	execCmd.Stdout = io.MultiWriter(os.Stdout, &outBuffer)
	execCmd.Stderr = io.MultiWriter(os.Stderr, &outBuffer)

	err := execCmd.Run()
	return outBuffer.Bytes(), err
}

func realLookPathFunc(file string) (string, error) {
	path, err := exec.LookPath(file)
	return path, err
}
