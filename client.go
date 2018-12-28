package vagrant_go

import (
	"github.com/palantir/stacktrace"
	"os/exec"
	"strings"
)

type Client struct {
	config         *Config
	commandRunFunc func(cmd string, args ...string) ([]byte, error)
	Box            BoxAPI
}

func NewClient(
	config *Config,
	commandRunFunc func(cmd string, args ...string) ([]byte, error),
) (*Client, error) {
	clientConfig := DefaultConfig()

	if config != nil {
		if len(config.BinaryName) > 0 {
			clientConfig.BinaryName = config.BinaryName
		}
	}

	_, err := exec.LookPath(clientConfig.BinaryName)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"`%s` not found in $PATH",
			clientConfig.BinaryName,
		)
	}

	clientCommandRunFunc := realCommandFunc
	if commandRunFunc != nil {
		clientCommandRunFunc = commandRunFunc
	}

	client := &Client{
		config:         clientConfig,
		commandRunFunc: clientCommandRunFunc,
	}

	client.Box = &boxAPI{
		client: client,
	}

	return client, nil
}

func (c *Client) executeVagrantCommand(args ...string) ([]*vagrantOutputLine, error) {
	cmdArgs := []string{
		"--machine-readable",
	}

	for _, arg := range args {
		cmdArgs = append(cmdArgs, arg)
	}

	output, err := c.commandRunFunc(c.config.BinaryName, cmdArgs...)
	return c.parseMachineReadableOutput(string(output)), err

}

func (c *Client) parseMachineReadableOutput(output string) []*vagrantOutputLine {
	var vagrantOutputLines []*vagrantOutputLine

	output = strings.Replace(output, `\n`, "\n", -1)
	outputLines := strings.Split(output, "\n")

	for _, outputLine := range outputLines {
		vagrantOutputLine := vagrantOutputLineFromString(outputLine)
		if vagrantOutputLine == nil {
			continue
		}

		vagrantOutputLines = append(
			vagrantOutputLines,
			vagrantOutputLine,
		)
	}

	return vagrantOutputLines
}
