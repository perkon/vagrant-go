package vagrant_go

import (
	"strings"
)

// Compile-time proof of interface implementation.
var _ GlobalAPI = (*globalAPI)(nil)

type GlobalAPI interface {
	Up(options *UpOptions) error
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
