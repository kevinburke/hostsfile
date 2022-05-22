//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package main

import (
	"os"
	"syscall"
)

// tempFile creates a new temporary hosts-file in an appropriate directory,
// opens the file for writing, and returns the resulting *os.File.
func tempFile(hostsPath string) (*os.File, error) {
	fs, err := os.Stat(hostsPath)
	if err != nil {
		return nil, err
	}
	f, err := os.CreateTemp("", "hostsfile-temp")
	if err != nil {
		return nil, err
	}

	// Set file mode to the same as the hosts-file.
	if err = os.Chmod(f.Name(), fs.Mode()); err != nil {
		return nil, err
	}

	// Set ownership to the same as the hosts-file.
	uid := fs.Sys().(*syscall.Stat_t).Uid
	gid := fs.Sys().(*syscall.Stat_t).Gid
	if err = os.Chown(f.Name(), int(uid), int(gid)); err != nil {
		return nil, err
	}

	return f, nil
}
