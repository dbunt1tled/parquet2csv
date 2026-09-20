package file

import (
	"errors"
	"os"
)

const FlushCount = 10000

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, err
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func IsWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, errors.New("Path doesn't exist. " + path)
	}

	if !info.IsDir() {
		return false, errors.New("Path isn't a directory. " + path)
	}

	return hasWriteAccess(path)
}
