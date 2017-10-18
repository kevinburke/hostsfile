package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Default location of hosts-file on Windows.
var hostsFile = os.Getenv("SystemRoot") + "\\System32\\drivers\\etc\\hosts"

// tempFile creates a new temporary hosts-file in an appropriate directory,
// opens the file for writing, and returns the resulting *os.File.
func tempFile(hostsPath string) (*os.File, error) {
	// Create the temporary file in the same location as the hosts-file to inherit
	// the correct permissions from the parent directory.
	return ioutil.TempFile(filepath.Dir(hostsPath), "hostsfile-temp")
}
