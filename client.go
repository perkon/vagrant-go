package vagrant_go

import (
	"github.com/palantir/stacktrace"
	"strings"
)

type Client struct {
	Config         *Config
	commandRunFunc func(cmd string, args ...string) ([]byte, error)
	osExecutor     OsExecutor
	Box            BoxAPI
	Global         GlobalAPI
}

func NewClient(
	config *Config,
	commandRunFunc func(cmd string, args ...string) ([]byte, error),
	lookPathFunc func(file string) (string, error),
) (*Client, error) {
	clientConfig := DefaultConfig()

	if config != nil && len(config.BinaryName) > 0 {
		clientConfig.BinaryName = config.BinaryName
	}

	clientLookPathFunc := realLookPathFunc
	if lookPathFunc != nil {
		clientLookPathFunc = lookPathFunc
	}

	_, err := clientLookPathFunc(clientConfig.BinaryName)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"`%s` not found in $PATH",
			clientConfig.BinaryName,
		)
	}

	clientCommandRunFunc := realCommandRunFunc
	if commandRunFunc != nil {
		clientCommandRunFunc = commandRunFunc
	}

	client := &Client{
		Config:         clientConfig,
		commandRunFunc: clientCommandRunFunc,
	}

	client.Box = &boxAPI{
		client: client,
	}

	client.Global = &globalAPI{
		client:     client,
		osExecutor: &osExecutor{},
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

	output, err := c.commandRunFunc(c.Config.BinaryName, cmdArgs...)
	return c.parseMachineReadableOutput(string(output)), err

}

func (c *Client) parseMachineReadableOutput(output string) []*vagrantOutputLine {
	vagrantOutputLines := []*vagrantOutputLine{}

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
