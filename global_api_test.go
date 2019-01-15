package vagrant_go

import (
	"errors"
	"github.com/stretchr/testify/assert"
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

func TestGlobalAPI_Up(t *testing.T) {
	t.Parallel()

	t.Run(
		"with default options and no execution error, it executes command and does not change current working dir before execution",
		func(t *testing.T) {
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
}
