package vagrant_go

import (
	"fmt"
	"github.com/kevinburke/ssh_config"
	"github.com/palantir/stacktrace"
	"strings"
)

// Compile-time proof of interface implementation.
var _ GlobalAPI = (*globalAPI)(nil)

type GlobalAPI interface {
	Up(options *UpOptions) error
	Destroy(options *DestroyOptions) error
	SshConfig(options *SshConfigOptions) (*ssh_config.Config, error)
}

type globalAPI struct {
	client     *Client
	osExecutor OsExecutor
}

type UpOptions struct {
	WorkingDirectory string
	Provision        bool
	ProvisionWith    []string
	DestroyOnError   bool
	Parallel         bool
	Provider         string
	InstallProvider  bool
}

func DefaultUpOptions() *UpOptions {
	return &UpOptions{
		WorkingDirectory: "",
		Provision:        true,
		ProvisionWith:    []string{},
		DestroyOnError:   true,
		Parallel:         true,
		Provider:         "",
		InstallProvider:  true,
	}
}

type DestroyOptions struct {
	WorkingDirectory string
	Force            bool
	Parallel         bool
}

func DefaultDestroyOptions() *DestroyOptions {
	return &DestroyOptions{
		WorkingDirectory: "",
		Force:            true,
		Parallel:         true,
	}
}

type SshConfigOptions struct {
	WorkingDirectory string
	Name             string
}

func DefaultSshConfigOptions() *SshConfigOptions {
	return &SshConfigOptions{
		WorkingDirectory: "",
		Name:             "",
	}
}

func (api *globalAPI) Up(options *UpOptions) error {
	args := []string{
		"up",
	}

	if options.Provision {
		args = append(args, "--provision")
	} else {
		args = append(args, "--no-provision")
	}

	if len(options.ProvisionWith) > 0 {
		args = append(args, "--provision-with", strings.Join(options.ProvisionWith, ","))
	}

	if options.DestroyOnError {
		args = append(args, "--destroy-on-error")
	} else {
		args = append(args, "--no-destroy-on-error")
	}

	if options.Parallel {
		args = append(args, "--parallel")
	} else {
		args = append(args, "--no-parallel")
	}

	if len(options.Provider) > 0 {
		args = append(args, "--provider", options.Provider)
	}

	if options.InstallProvider {
		args = append(args, "--install-provider")
	} else {
		args = append(args, "--no-install-provider")
	}

	if len(options.WorkingDirectory) > 0 {
		oldWorkingDir, err := api.osExecutor.Getwd()
		if err != nil {
			return err
		}

		err = api.osExecutor.Chdir(options.WorkingDirectory)
		if err != nil {
			return err
		}

		_, err = api.client.executeVagrantCommand(args...)
		if err != nil {
			return err
		}

		err = api.osExecutor.Chdir(oldWorkingDir)
		return err
	}

	_, err := api.client.executeVagrantCommand(args...)
	return err
}

func (api *globalAPI) Destroy(options *DestroyOptions) error {
	args := []string{
		"destroy",
	}

	if options.Force {
		args = append(args, "--force")
	}

	if options.Parallel {
		args = append(args, "--parallel")
	} else {
		args = append(args, "--no-parallel")
	}

	if len(options.WorkingDirectory) > 0 {
		oldWorkingDir, err := api.osExecutor.Getwd()
		if err != nil {
			return err
		}

		err = api.osExecutor.Chdir(options.WorkingDirectory)
		if err != nil {
			return err
		}

		_, err = api.client.executeVagrantCommand(args...)
		if err != nil {
			return err
		}

		err = api.osExecutor.Chdir(oldWorkingDir)
		return err
	}

	_, err := api.client.executeVagrantCommand(args...)
	return err
}

func (api *globalAPI) SshConfig(options *SshConfigOptions) (*ssh_config.Config, error) {
	args := []string{
		"ssh-config",
	}

	if len(options.Name) > 0 {
		args = append(args, "--name", options.Name)
	}

	var oldWorkingDir string
	var err error

	hasWorkingDir := len(options.WorkingDirectory) > 0
	if hasWorkingDir {
		oldWorkingDir, err = api.osExecutor.Getwd()
		if err != nil {
			return nil, err
		}

		err = api.osExecutor.Chdir(options.WorkingDirectory)
		if err != nil {
			return nil, err
		}
	}

	outputLines, err := api.client.executeVagrantCommand(args...)
	if err != nil {
		return nil, err
	}

	var sshConfig *ssh_config.Config
	var output string

	for _, line := range outputLines {
		if len(line.data) < 1 || line.kind != "ssh-config" {
			continue
		}

		hostConfig := strings.Replace(line.data[0], `\n`, "\n", -1)
		output += fmt.Sprintf("%s\n", hostConfig)
	}

	sshConfig, err = ssh_config.Decode(
		strings.NewReader(output),
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to decode ssh_config")
	}

	if hasWorkingDir {
		err = api.osExecutor.Chdir(oldWorkingDir)
		if err != nil {
			return nil, err
		}
	}

	return sshConfig, nil
}
