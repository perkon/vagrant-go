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
