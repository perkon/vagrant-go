package vagrant_go

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run(
		"with nil `config`, `commandRunFunc` and `lookPathFunc` given, it uses default config and probably `realCommandRunFunc` and `realLookPathFunc`",
		func(t *testing.T) {
			client, err := NewClient(nil, nil, nil)
			require.NoError(t, err)

			assert.NotNil(t, client)
			assert.Equal(t, defaultBinaryName, client.Config.BinaryName)
			assert.NotNil(t, client.commandRunFunc)

			assert.NotNil(t, client.Box)
		},
	)
}

func TestParseMachineReadableOutput(t *testing.T) {
	t.Parallel()

	t.Run(
		"with machine readable output in format 'timestamp,target,type,data...' as per https://www.vagrantup.com/docs/cli/machine-readable.html, it returns slice of parsed lines",
		func(t *testing.T) {
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
				`The Libvirt domain is running. To stop this machine%!(VAGRANT_COMMA) you can run`,
				actualLine.data[0],
			)
		},
	)

	t.Run(
		"with blank output, it returns empty slice of lines",
		func(t *testing.T) {
			client := emptyTestClient(t)

			lines := client.parseMachineReadableOutput("")
			require.NotNil(t, lines)
			assert.Empty(t, lines)
		},
	)

	t.Run(
		"with non-machine readable output, that is not blank, it returns empty slice of lines",
		func(t *testing.T) {
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
			client := emptyTestClient(t)

			lines := client.parseMachineReadableOutput("")
			require.NotNil(t, lines)
			assert.Empty(t, lines)
		},
	)
}
