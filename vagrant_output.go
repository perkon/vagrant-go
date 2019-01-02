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
	data      []string
}

func vagrantOutputLineFromString(str string) *vagrantOutputLine {
	trimmedStr := strings.TrimSpace(str)
	splitLines := strings.SplitN(trimmedStr, ",", 4)

	//noinspection GoPreferNilSlice
	vagrantLines := []string{}

	for _, line := range splitLines {
		if contains(ignoredOutputLines, line) {
			continue
		}

		vagrantLines = append(vagrantLines, line)
	}

	if len(vagrantLines) < 4 {
		return nil
	}

	//noinspection GoPreferNilSlice
	dataLines := []string{}

	for _, line := range strings.Split(vagrantLines[3], ",") {
		dataLines = append(dataLines, line)
	}

	return &vagrantOutputLine{
		timestamp: vagrantLines[0],
		target:    vagrantLines[1],
		kind:      vagrantLines[2],
		data:      dataLines,
	}
}
