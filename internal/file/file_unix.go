//go:build linux || darwin

package file

import (
	"fmt"
	"syscall"
)

// writeOK is POSIX W_OK, which the standard syscall package does not export.
const writeOK = 0x2

// hasWriteAccess asks the kernel instead of reading mode bits, so group and other
// permissions count too, and a read-only mount is reported as such.
func hasWriteAccess(path string) (bool, error) {
	if err := syscall.Access(path, writeOK); err != nil {
		return false, fmt.Errorf("no write access to %s: %w", path, err)
	}
	return true, nil
}
