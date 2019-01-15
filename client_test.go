package vagrant_go

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run(
		"with nil `config`, `commandRunFunc` and `lookPathFunc` given, it uses default config and probably `realCommandRunFunc` and `realLookPathFunc`",
		func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(nil, nil, nil)
			require.NoError(t, err)

			assert.NotNil(t, client)
			assert.Equal(t, defaultBinaryName, client.Config.BinaryName)
			assert.NotNil(t, client.commandRunFunc)

			assert.NotNil(t, client.Box)
			assert.NotNil(t, client.Global)
		},
	)

	t.Run(
		"with 'config.BinaryName' = vagrant123, it uses provided 'BinaryName'",
		func(t *testing.T) {
			t.Parallel()

			config := DefaultConfig()
			config.BinaryName = "vagrant123"

			lookPathFunc := func(binaryName string) (string, error) {
				return "/tmp/vagrant123", nil
			}

			client, err := NewClient(config, nil, lookPathFunc)
			require.NoError(t, err)

			assert.NotNil(t, client)
			assert.Equal(t, "vagrant123", client.Config.BinaryName)
			assert.NotNil(t, client.commandRunFunc)

			assert.NotNil(t, client.Box)
			assert.NotNil(t, client.Global)
		},
	)

	t.Run(
		"with 'clientLookPathFunc' returning an error, it returns error that 'config.BinaryName' is not found",
		func(t *testing.T) {
			t.Parallel()

			config := DefaultConfig()
			config.BinaryName = "vagrant123"

			isLookPathFuncCalled := false

			clientLookPathFunc := func(binaryName string) (string, error) {
				assert.Equal(t, "vagrant123", binaryName)
				isLookPathFuncCalled = true
				return "/tmp/example", errors.New("fake error")
			}

			client, err := NewClient(
				config,
				nil,
				clientLookPathFunc,
			)
			assert.Error(t, err, "fake error")

			assert.Nil(t, client)
			assert.True(t, isLookPathFuncCalled)
		},
	)
}

func TestExecuteVagrantCommand(t *testing.T) {
	t.Run(
		"when there's no command execution error, it executes command with given 'args' and returns parsed machine readable output and no error",
		func(t *testing.T) {
			t.Parallel()

			client := emptyTestClient(t)
			isCommandRunCalled := false
			args := []string{"version"}

			fakeOutput := "1546430404,default,provider-name,libvirt"
			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 2)
				isCommandRunCalled = true
				return []byte(fakeOutput), nil
			}

			outputLines, err := client.executeVagrantCommand(args...)
			require.NoError(t, err)

			assert.Equal(t, len(outputLines), 1)
			assert.Equal(t, outputLines[0].timestamp, "1546430404")
			assert.Equal(t, outputLines[0].target, "default")
			assert.Equal(t, outputLines[0].kind, "provider-name")
			assert.Equal(t, len(outputLines[0].data), 1)
			assert.Equal(t, outputLines[0].data[0], "libvirt")

			assert.True(t, isCommandRunCalled)
		},
	)

	t.Run(
		"when there's a command execution error, it executes command with given 'args' and returns parsed machine readable output and error",
		func(t *testing.T) {
			t.Parallel()

			client := emptyTestClient(t)
			isCommandRunCalled := false
			args := []string{"version"}

			fakeOutput := "1546430404,default,provider-name,libvirt"
			fakeErrorMessage := "fakeCommandRunError"

			client.commandRunFunc = func(cmd string, args ...string) (bytes []byte, e error) {
				assert.Equal(t, cmd, client.Config.BinaryName)
				assert.Len(t, args, 2)
				assert.Equal(t, args[0], "--machine-readable")
				assert.Equal(t, args[1], "version")
				isCommandRunCalled = true
				return []byte(fakeOutput), errors.New(fakeErrorMessage)
			}

			outputLines, err := client.executeVagrantCommand(args...)
			require.Error(t, err, fakeErrorMessage)

			assert.Equal(t, len(outputLines), 1)
			assert.Equal(t, outputLines[0].timestamp, "1546430404")
			assert.Equal(t, outputLines[0].target, "default")
			assert.Equal(t, outputLines[0].kind, "provider-name")
			assert.Equal(t, len(outputLines[0].data), 1)
			assert.Equal(t, outputLines[0].data[0], "libvirt")

			assert.True(t, isCommandRunCalled)
		},
	)
}

func TestParseMachineReadableOutput(t *testing.T) {
	t.Run(
		"with machine readable output in format 'timestamp,target,type,data...' as per https://www.vagrantup.com/docs/cli/machine-readable.html, it returns slice of parsed lines",
		func(t *testing.T) {
			t.Parallel()

			client := emptyTestClient(t)

			output := `
1546430404,default,metadata,provider,libvirt
1546430404,default,provider-name,libvirt
1546430404,default,state,running
1546430404,default,state-human-short,running
1546430404,default,state-human-long,The Libvirt domain is running. To stop this machine%!(VAGRANT_COMMA) you can run\n'vagrant halt''. To destroy the machine%!(VAGRANT_COMMA) you can run 'vagrant destroy'.
1546430404,,ui,info,Current machine states:\n\ndefault                   running (libvirt)\n\nThe Libvirt domain is running. To stop this machine%!(VAGRANT_COMMA) you can run\n'vagrant halt''. To destroy the machine%!(VAGRANT_COMMA) you can run 'vagrant destroy'.
`

			lines := client.parseMachineReadableOutput(output)
			require.NotNil(t, lines)
			require.Len(t, lines, 4)

			actualLine := lines[0]
			assert.Equal(t, "1546430404", actualLine.timestamp)
			assert.Equal(t, "default", actualLine.target)
			assert.Equal(t, "provider-name", actualLine.kind)
			require.Len(t, actualLine.data, 1)
			assert.Equal(t, "libvirt", actualLine.data[0])

			actualLine = lines[1]
			assert.Equal(t, "1546430404", actualLine.timestamp)
			assert.Equal(t, "default", actualLine.target)
			assert.Equal(t, "state", actualLine.kind)
			require.Len(t, actualLine.data, 1)
			assert.Equal(t, "running", actualLine.data[0])

			actualLine = lines[2]
			assert.Equal(t, "1546430404", actualLine.timestamp)
			assert.Equal(t, "default", actualLine.target)
			assert.Equal(t, "state-human-short", actualLine.kind)
			require.Len(t, actualLine.data, 1)
			assert.Equal(t, "running", actualLine.data[0])

			actualLine = lines[3]
			assert.Equal(t, "1546430404", actualLine.timestamp)
			assert.Equal(t, "default", actualLine.target)
			assert.Equal(t, "state-human-long", actualLine.kind)
			require.Len(t, actualLine.data, 1)
			// NOTE: It's expected to lose the full actual output due to newlines
			assert.Equal(t,
				`The Libvirt domain is running. To stop this machine%!(VAGRANT_COMMA) you can run\n'vagrant halt''. To destroy the machine%!(VAGRANT_COMMA) you can run 'vagrant destroy'.`,
				actualLine.data[0],
			)
		},
	)

	t.Run(
		"with blank output, it returns empty slice of lines",
		func(t *testing.T) {
			t.Parallel()
			client := emptyTestClient(t)

			lines := client.parseMachineReadableOutput("")
			require.NotNil(t, lines)
			assert.Empty(t, lines)
		},
	)

	t.Run(
		"with non-machine readable output, that is not blank, it returns empty slice of lines",
		func(t *testing.T) {
			t.Parallel()
			client := emptyTestClient(t)

			// NOTE: Can you think of a more ridiculous output? :)
			output := `
(Not all processes could be identified, non-owned process info                                                                                                                                                                                                                  
 will not be shown, you would have to be root to see it all.)                                                                                                                                                                                                                   
Active Internet connections (servers and established)                                                                                                                                                                                                                           
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name                                                                                                                                                                                
tcp        0      0 127.0.0.1:63342         0.0.0.0:*               LISTEN      7624/java                                                                                                                                                                                       
tcp        0      0 0.0.0.0:111             0.0.0.0:*               LISTEN      -                                                                                                                                                                                               
tcp        0      0 0.0.0.0:57621           0.0.0.0:*               LISTEN      7376/spotify                                                                                                                                                                                    
tcp        0      0 0.0.0.0:39351           0.0.0.0:*               LISTEN      7376/spotify                                                                                                                                                                                    
tcp        0      0 127.0.0.1:631           0.0.0.0:*               LISTEN      -                                                                                                                                                                                               
tcp        0      0 127.0.0.1:6942          0.0.0.0:*               LISTEN      7624/java                                                                                                                                                                                       
tcp        0      0 127.0.0.1:38183         0.0.0.0:*               LISTEN      4958/usr/share/jetb                                                                                                                                                                             
tcp        0      0 192.168.13.37:43778     192.121.140.177:80      ESTABLISHED 7376/spotify              
`

			lines := client.parseMachineReadableOutput(output)
			require.NotNil(t, lines)
			assert.Empty(t, lines)
		},
	)

	t.Run(
		"with blank output, it returns empty slice of lines",
		func(t *testing.T) {
			t.Parallel()
			client := emptyTestClient(t)

			lines := client.parseMachineReadableOutput("")
			require.NotNil(t, lines)
			assert.Empty(t, lines)
		},
	)
}
