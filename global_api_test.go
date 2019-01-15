package vagrant_go

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultUpOptions(t *testing.T) {
	t.Parallel()

	options := DefaultUpOptions()

	assert.Equal(t, options.WorkingDirectory, "")
	assert.True(t, options.Provision)
	assert.Equal(t, len(options.ProvisionWith), 0)
	assert.True(t, options.DestroyOnError)
	assert.True(t, options.Parallel)
	assert.Equal(t, options.Provider, "")
	assert.True(t, options.InstallProvider)
}

func TestDefaultDestroyOptions(t *testing.T) {
	t.Parallel()

	options := DefaultDestroyOptions()

	assert.Equal(t, options.WorkingDirectory, "")
	assert.True(t, options.Force)
	assert.True(t, options.Parallel)
}

func TestDefaultSshConfigOptions(t *testing.T) {
	t.Parallel()

	options := DefaultSshConfigOptions()

	assert.Equal(t, options.WorkingDirectory, "")
	assert.Equal(t, options.Name, "")
}

func TestGlobalAPI_Up(t *testing.T) {
	t.Run(
		"with default options and no execution error, it executes command and does not change current working dir before execution",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			err := client.Global.Up(options)
			require.NoError(t, err)

			assert.True(t, isCommandRunCalled)

			fakeOsExecutor.AssertNotCalled(t, "Getwd")
			fakeOsExecutor.AssertNotCalled(t, "Chdir")
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when retrieving current working dir, it does not execute command and does not change current working dir before execution",
		func(t *testing.T) {
			t.Parallel()
			fakeError := errors.New("fake error")
			fakeOsExecutor := &fakeOsExecutor{}
			fakeOsExecutor.On("Getwd").Return("", fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.WorkingDirectory = "/tmp/example"
			err := client.Global.Up(options)
			assert.Error(t, err, "fake error")

			assert.False(t, isCommandRunCalled)
			fakeOsExecutor.AssertNumberOfCalls(t, "Getwd", 1)
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when changing current working dir to specified one in options 'workingDir', it does not execute command and does not change current working dir",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			fakeCwd := "/tmp/anotherexample"
			fakeOsExecutor.On("Getwd").Return(fakeCwd, nil)

			fakeOptionsWd := "/tmp/example"
			fakeError := errors.New("fake error")
			fakeOsExecutor.On("Chdir", fakeOptionsWd).Return(fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.WorkingDirectory = fakeOptionsWd

			err := client.Global.Up(options)
			assert.Error(t, err, "fake error")

			assert.False(t, isCommandRunCalled)
			fakeOsExecutor.AssertCalled(t, "Getwd")
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeOptionsWd)
			fakeOsExecutor.AssertNotCalled(t, "Chdir", fakeCwd)
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when changing current working dir to old one, it executes command and does not change current working dir to old one",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			fakeCwd := "/tmp/anotherexample"
			fakeOsExecutor.On("Getwd").Return(fakeCwd, nil)

			fakeOptionsWd := "/tmp/example"
			fakeError := errors.New("fake error")
			fakeOsExecutor.On("Chdir", fakeOptionsWd).Return(nil)
			fakeOsExecutor.On("Chdir", fakeCwd).Return(fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.WorkingDirectory = fakeOptionsWd

			err := client.Global.Up(options)
			assert.Error(t, err, "fake error")

			assert.True(t, isCommandRunCalled)
			fakeOsExecutor.AssertCalled(t, "Getwd")
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeOptionsWd)
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeCwd)
		},
	)

	t.Run(
		"with options providing 'Provision', it executes command with '--provision'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.Provision = true

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Provision' = false, it executes command with '--no-provision'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--no-provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.Provision = false

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'ProvisionWith' = [shell], it executes command with '--provision-with shell'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 8)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--provision-with")
				assert.Equal(t, args[4], "shell")
				assert.Equal(t, args[5], "--destroy-on-error")
				assert.Equal(t, args[6], "--parallel")
				assert.Equal(t, args[7], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.ProvisionWith = []string{"shell"}

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'DestroyOnError' = true, it executes command with '--destroy-on-error'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.DestroyOnError = true

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'DestroyOnError' = false, it executes command with '--no-destroy-on-error'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--no-destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.DestroyOnError = false

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Parallel' = true, it executes command with '--parallel'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.Parallel = true

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Parallel' = false, it executes command with '--no-parallel'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--no-parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.Parallel = false

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Provider' = 'libvirt', it executes command with '--provider libvirt'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 8)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--provider")
				assert.Equal(t, args[6], "libvirt")
				assert.Equal(t, args[7], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.Provider = "libvirt"

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Provider' = '', it executes command with no '--provider'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.Provider = ""

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'InstallProvider' = true, it executes command with no '--install-provider'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.InstallProvider = true

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run("with options providing 'WorkingDirectory' = /tmp/example' and command execution returning an error, it returns an error", func(t *testing.T) {
		t.Parallel()

		fakeOsExecutor := &fakeOsExecutor{}
		fakeOsExecutor.On("Getwd").Return("/tmp/example", nil)
		fakeOsExecutor.On("Chdir", mock.Anything).Return(nil)

		isCommandRunCalled := false
		client := emptyTestClient(t)
		client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
			isCommandRunCalled = true
			return []byte{}, errors.New("fake error")
		}

		globalAPI := &globalAPI{
			osExecutor: fakeOsExecutor,
			client:     client,
		}
		client.Global = globalAPI

		options := DefaultUpOptions()
		options.WorkingDirectory = "/tmp/example"
		err := client.Global.Up(options)
		assert.Error(t, err, "fake error")

		assert.True(t, isCommandRunCalled)
	})

	t.Run(
		"with options providing 'InstallProvider' = false, it executes command with no '--install-provider'",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 6)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "up")
				assert.Equal(t, args[2], "--provision")
				assert.Equal(t, args[3], "--destroy-on-error")
				assert.Equal(t, args[4], "--parallel")
				assert.Equal(t, args[5], "--no-install-provider")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.InstallProvider = false

			err := client.Global.Up(options)
			assert.NoError(t, err)

			assert.True(t, isCommandRunCalled)
		},
	)
}

func TestGlobalAPI_Destroy(t *testing.T) {
	t.Run(
		"with options providing 'Force' = true, it executes command with '--force'",
		func(t *testing.T) {
			t.Parallel()

			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI

			isCommandRunCalled := false
			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 4)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "destroy")
				assert.Equal(t, args[2], "--force")
				assert.Equal(t, args[3], "--parallel")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultDestroyOptions()
			options.Force = true
			err := client.Global.Destroy(options)
			require.NoError(t, err)
			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Force' = false, it executes command without '--force'",
		func(t *testing.T) {
			t.Parallel()

			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI

			isCommandRunCalled := false
			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 3)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "destroy")
				assert.Equal(t, args[2], "--parallel")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultDestroyOptions()
			options.Force = false
			err := client.Global.Destroy(options)
			require.NoError(t, err)
			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Parallel' = true, it executes command with '--parallel'",
		func(t *testing.T) {
			t.Parallel()

			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI

			isCommandRunCalled := false
			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 4)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "destroy")
				assert.Equal(t, args[2], "--force")
				assert.Equal(t, args[3], "--parallel")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultDestroyOptions()
			options.Parallel = true
			err := client.Global.Destroy(options)
			require.NoError(t, err)
			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with options providing 'Parallel' = false, it executes command with '--no-parallel'",
		func(t *testing.T) {
			t.Parallel()

			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI

			isCommandRunCalled := false
			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 4)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "destroy")
				assert.Equal(t, args[2], "--force")
				assert.Equal(t, args[3], "--no-parallel")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultDestroyOptions()
			options.Parallel = false
			err := client.Global.Destroy(options)
			require.NoError(t, err)
			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"with default options and no execution error, it executes command and does not change current working dir before execution",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 4)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "destroy")
				assert.Equal(t, args[2], "--force")
				assert.Equal(t, args[3], "--parallel")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultDestroyOptions()
			err := client.Global.Destroy(options)
			require.NoError(t, err)

			assert.True(t, isCommandRunCalled)

			fakeOsExecutor.AssertNotCalled(t, "Getwd")
			fakeOsExecutor.AssertNotCalled(t, "Chdir")
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when retrieving current working dir, it does not execute command and does not change current working dir before execution",
		func(t *testing.T) {
			t.Parallel()
			fakeError := errors.New("fake error")
			fakeOsExecutor := &fakeOsExecutor{}
			fakeOsExecutor.On("Getwd").Return("", fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultUpOptions()
			options.WorkingDirectory = "/tmp/example"
			err := client.Global.Up(options)
			assert.Error(t, err, "fake error")

			assert.False(t, isCommandRunCalled)
			fakeOsExecutor.AssertNumberOfCalls(t, "Getwd", 1)
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when changing current working dir to specified one in options 'workingDir', it does not execute command and does not change current working dir",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			fakeCwd := "/tmp/anotherexample"
			fakeOsExecutor.On("Getwd").Return(fakeCwd, nil)

			fakeOptionsWd := "/tmp/example"
			fakeError := errors.New("fake error")
			fakeOsExecutor.On("Chdir", fakeOptionsWd).Return(fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultDestroyOptions()
			options.WorkingDirectory = fakeOptionsWd

			err := client.Global.Destroy(options)
			assert.Error(t, err, "fake error")

			assert.False(t, isCommandRunCalled)
			fakeOsExecutor.AssertCalled(t, "Getwd")
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeOptionsWd)
			fakeOsExecutor.AssertNotCalled(t, "Chdir", fakeCwd)
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when changing current working dir to old one, it executes command and does not change current working dir to old one",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			fakeCwd := "/tmp/anotherexample"
			fakeOsExecutor.On("Getwd").Return(fakeCwd, nil)

			fakeOptionsWd := "/tmp/example"
			fakeError := errors.New("fake error")
			fakeOsExecutor.On("Chdir", fakeOptionsWd).Return(nil)
			fakeOsExecutor.On("Chdir", fakeCwd).Return(fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultDestroyOptions()
			options.WorkingDirectory = fakeOptionsWd

			err := client.Global.Destroy(options)
			assert.Error(t, err, "fake error")

			assert.True(t, isCommandRunCalled)
			fakeOsExecutor.AssertCalled(t, "Getwd")
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeOptionsWd)
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeCwd)
		},
	)
}

func TestGlobalAPI_SshConfig(t *testing.T) {
	t.Run(
		"with default options and no execution error, it executes command and does not change current working dir before execution",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 2)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "ssh-config")

				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultSshConfigOptions()
			_, err := client.Global.SshConfig(options)
			require.NoError(t, err)

			assert.True(t, isCommandRunCalled)

			fakeOsExecutor.AssertNotCalled(t, "Getwd")
			fakeOsExecutor.AssertNotCalled(t, "Chdir")
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when retrieving current working dir, it does not execute command and does not change current working dir before execution",
		func(t *testing.T) {
			t.Parallel()
			fakeError := errors.New("fake error")
			fakeOsExecutor := &fakeOsExecutor{}
			fakeOsExecutor.On("Getwd").Return("", fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultSshConfigOptions()
			options.WorkingDirectory = "/tmp/example"
			_, err := client.Global.SshConfig(options)
			assert.Error(t, err, "fake error")

			assert.False(t, isCommandRunCalled)
			fakeOsExecutor.AssertNumberOfCalls(t, "Getwd", 1)
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when changing current working dir to specified one in options 'workingDir', it does not execute command and does not change current working dir",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			fakeCwd := "/tmp/anotherexample"
			fakeOsExecutor.On("Getwd").Return(fakeCwd, nil)

			fakeOptionsWd := "/tmp/example"
			fakeError := errors.New("fake error")
			fakeOsExecutor.On("Chdir", fakeOptionsWd).Return(fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultSshConfigOptions()
			options.WorkingDirectory = fakeOptionsWd

			_, err := client.Global.SshConfig(options)
			assert.Error(t, err, "fake error")

			assert.False(t, isCommandRunCalled)
			fakeOsExecutor.AssertCalled(t, "Getwd")
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeOptionsWd)
			fakeOsExecutor.AssertNotCalled(t, "Chdir", fakeCwd)
		},
	)

	t.Run(
		"with options providing 'workingDir' and an error when changing current working dir to old one, it executes command and does not change current working dir to old one",
		func(t *testing.T) {
			t.Parallel()
			fakeOsExecutor := &fakeOsExecutor{}

			fakeCwd := "/tmp/anotherexample"
			fakeOsExecutor.On("Getwd").Return(fakeCwd, nil)

			fakeOptionsWd := "/tmp/example"
			fakeError := errors.New("fake error")
			fakeOsExecutor.On("Chdir", fakeOptionsWd).Return(nil)
			fakeOsExecutor.On("Chdir", fakeCwd).Return(fakeError)

			client := emptyTestClient(t)
			globalAPI := &globalAPI{
				osExecutor: fakeOsExecutor,
				client:     client,
			}
			client.Global = globalAPI
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true
				return []byte{}, nil
			}

			options := DefaultSshConfigOptions()
			options.WorkingDirectory = fakeOptionsWd

			_, err := client.Global.SshConfig(options)
			assert.Error(t, err, "fake error")

			assert.True(t, isCommandRunCalled)
			fakeOsExecutor.AssertCalled(t, "Getwd")
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeOptionsWd)
			fakeOsExecutor.AssertCalled(t, "Chdir", fakeCwd)
		},
	)

	t.Run(
		"with single-host SSH config returned, it returns SshConfig for a single host",
		func(t *testing.T) {
			t.Parallel()

			client := emptyTestClient(t)
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true

				output := `
1547581456,default,metadata,provider,libvirt
1547581456,default,ssh-config,Host default\n  HostName 192.168.121.165\n  User root\n  Port 22\n  UserKnownHostsFile /dev/null\n  StrictHostKeyChecking no\n  PasswordAuthentication no\n  IdentityFile /tmp/example123456789/.vagrant/machines/default/libvirt/private_key\n  IdentitiesOnly yes\n  LogLevel FATAL\n
Host default
  HostName 192.168.121.165
  User root
  Port 22
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /tmp/example123456789/.vagrant/machines/default/libvirt/private_key
  IdentitiesOnly yes
  LogLevel FATAL
`
				return []byte(output), nil
			}

			options := DefaultSshConfigOptions()
			sshConfig, err := client.Global.SshConfig(options)
			require.NoError(t, err)
			assert.True(t, isCommandRunCalled)
			require.NotNil(t, sshConfig)

			hostname, err := sshConfig.Get("default", "HostName")
			require.NoError(t, err)
			assert.Equal(t, hostname, "192.168.121.165")

			user, err := sshConfig.Get("default", "User")
			require.NoError(t, err)
			assert.Equal(t, user, "root")

			port, err := sshConfig.Get("default", "Port")
			require.NoError(t, err)
			assert.Equal(t, port, "22")

			userKnownHostsFile, err := sshConfig.Get("default", "UserKnownHostsFile")
			require.NoError(t, err)
			assert.Equal(t, userKnownHostsFile, "/dev/null")

			strictHostKeyChecking, err := sshConfig.Get("default", "StrictHostKeyChecking")
			require.NoError(t, err)
			assert.Equal(t, strictHostKeyChecking, "no")

			passwordAuthentication, err := sshConfig.Get("default", "PasswordAuthentication")
			require.NoError(t, err)
			assert.Equal(t, passwordAuthentication, "no")

			identityFile, err := sshConfig.Get("default", "IdentityFile")
			require.NoError(t, err)
			assert.Equal(
				t,
				identityFile,
				"/tmp/example123456789/.vagrant/machines/default/libvirt/private_key",
			)

			identitiesOnly, err := sshConfig.Get("default", "IdentitiesOnly")
			require.NoError(t, err)
			assert.Equal(
				t,
				identitiesOnly,
				"yes",
			)

			logLevel, err := sshConfig.Get("default", "LogLevel")
			require.NoError(t, err)
			assert.Equal(
				t,
				logLevel,
				"FATAL",
			)
		},
	)

	t.Run(
		"with 3 hosts SSH config returned, it returns SshConfig for the 3 hosts",
		func(t *testing.T) {
			t.Parallel()

			client := emptyTestClient(t)
			isCommandRunCalled := false

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				isCommandRunCalled = true

				output := `
1547587389,master,metadata,provider,libvirt
1547587389,node1,metadata,provider,libvirt
1547587389,node2,metadata,provider,libvirt
1547587390,master,ssh-config,Host master\n  HostName 192.168.121.148\n  User vagrant\n  Port 22\n  UserKnownHostsFile /dev/null\n  StrictHostKeyChecking no\n  PasswordAuthentication no\n  IdentityFile /home/syndbg/.vagrant.d/insecure_private_key\n  IdentitiesOnly yes\n  LogLevel FATAL\n
Host master
  HostName 192.168.121.148
  User vagrant
  Port 22
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /home/syndbg/.vagrant.d/insecure_private_key
  IdentitiesOnly yes
  LogLevel FATAL

1547587390,node1,ssh-config,Host node1\n  HostName 192.168.121.223\n  User vagrant\n  Port 22\n  UserKnownHostsFile /dev/null\n  StrictHostKeyChecking no\n  PasswordAuthentication no\n  IdentityFile /home/syndbg/.vagrant.d/insecure_private_key\n  IdentitiesOnly yes\n  LogLevel FATAL\n
Host node1
  HostName 192.168.121.223
  User vagrant
  Port 22
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /home/syndbg/.vagrant.d/insecure_private_key
  IdentitiesOnly yes
  LogLevel FATAL

1547587390,node2,ssh-config,Host node2\n  HostName 192.168.121.206\n  User vagrant\n  Port 22\n  UserKnownHostsFile /dev/null\n  StrictHostKeyChecking no\n  PasswordAuthentication no\n  IdentityFile /home/syndbg/.vagrant.d/insecure_private_key\n  IdentitiesOnly yes\n  LogLevel FATAL\n
Host node2
  HostName 192.168.121.206
  User vagrant
  Port 22
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /home/syndbg/.vagrant.d/insecure_private_key
  IdentitiesOnly yes
  LogLevel FATAL
`
				return []byte(output), nil
			}

			options := DefaultSshConfigOptions()
			sshConfig, err := client.Global.SshConfig(options)
			require.NoError(t, err)
			assert.True(t, isCommandRunCalled)
			require.NotNil(t, sshConfig)

			tests := []struct {
				Host                   string
				HostName               string
				User                   string
				Port                   string
				UserKnownHostsFile     string
				StrictHostKeyChecking  string
				PasswordAuthentication string
				IdentityFile           string
				IdentitiesOnly         string
				LogLevel               string
			}{
				{
					Host:                   "master",
					HostName:               "192.168.121.148",
					User:                   "vagrant",
					Port:                   "22",
					UserKnownHostsFile:     "/dev/null",
					StrictHostKeyChecking:  "no",
					PasswordAuthentication: "no",
					IdentityFile:           "/home/syndbg/.vagrant.d/insecure_private_key",
					IdentitiesOnly:         "yes",
					LogLevel:               "FATAL",
				},
				{
					Host:                   "node1",
					HostName:               "192.168.121.223",
					User:                   "vagrant",
					Port:                   "22",
					UserKnownHostsFile:     "/dev/null",
					StrictHostKeyChecking:  "no",
					PasswordAuthentication: "no",
					IdentityFile:           "/home/syndbg/.vagrant.d/insecure_private_key",
					IdentitiesOnly:         "yes",
					LogLevel:               "FATAL",
				},
				{
					Host:                   "node2",
					HostName:               "192.168.121.206",
					User:                   "vagrant",
					Port:                   "22",
					UserKnownHostsFile:     "/dev/null",
					StrictHostKeyChecking:  "no",
					PasswordAuthentication: "no",
					IdentityFile:           "/home/syndbg/.vagrant.d/insecure_private_key",
					IdentitiesOnly:         "yes",
					LogLevel:               "FATAL",
				},
			}

			for _, subTest := range tests {
				hostname, err := sshConfig.Get(subTest.Host, "HostName")
				require.NoError(t, err)
				assert.Equal(t, hostname, subTest.HostName)

				user, err := sshConfig.Get(subTest.Host, "User")
				require.NoError(t, err)
				assert.Equal(t, user, subTest.User)

				port, err := sshConfig.Get(subTest.Host, "Port")
				require.NoError(t, err)
				assert.Equal(t, port, subTest.Port)

				userKnownHostsFile, err := sshConfig.Get(subTest.Host, "UserKnownHostsFile")
				require.NoError(t, err)
				assert.Equal(t, userKnownHostsFile, subTest.UserKnownHostsFile)

				strictHostKeyChecking, err := sshConfig.Get(
					subTest.Host,
					"StrictHostKeyChecking",
				)
				require.NoError(t, err)
				assert.Equal(t, strictHostKeyChecking, subTest.StrictHostKeyChecking)

				passwordAuthentication, err := sshConfig.Get(subTest.Host, "PasswordAuthentication")
				require.NoError(t, err)
				assert.Equal(t, passwordAuthentication, subTest.PasswordAuthentication)

				identityFile, err := sshConfig.Get(subTest.Host, "IdentityFile")
				require.NoError(t, err)
				assert.Equal(
					t,
					identityFile,
					subTest.IdentityFile,
				)

				identitiesOnly, err := sshConfig.Get(subTest.Host, "IdentitiesOnly")
				require.NoError(t, err)
				assert.Equal(
					t,
					identitiesOnly,
					subTest.IdentitiesOnly,
				)

				logLevel, err := sshConfig.Get(subTest.Host, "LogLevel")
				require.NoError(t, err)
				assert.Equal(
					t,
					logLevel,
					subTest.LogLevel,
				)
			}
		},
	)
}
