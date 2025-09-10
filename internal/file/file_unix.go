//go:build linux || darwin

package file

import (
	"errors"
	"os"
	"syscall"
)

func isWritableByOwner(path string) (bool, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return false, errors.New("Unable to get stat. " + path)
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, errors.New("User doesn't have permission to write to this directory. " + path)
	}
	return true, nil
}