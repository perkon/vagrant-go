package vagrant_go

import (
	"strings"
)

var ignoredOutputLines = []string{"metadata", "ui", "action"}

// NOTE: A Vagrant line is in format of `timestamp,target,type,data`.
// ref: https://www.vagrantup.com/docs/cli/machine-readable.html
type vagrantOutputLine struct {
	timestamp string
	target    string
	kind      string
	data      string
}

func vagrantOutputLineFromString(str string) *vagrantOutputLine {
	vagrantLines := strings.SplitN(str, ",", 4)

	if len(vagrantLines) < 4 || contains(ignoredOutputLines, vagrantLines[1]) {
		return nil
	}

	return &vagrantOutputLine{
		timestamp: vagrantLines[0],
		target:    vagrantLines[1],
		kind:      vagrantLines[2],
		data:      vagrantLines[3],
	}
}
