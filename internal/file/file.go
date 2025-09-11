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
		panic(err)
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

	if info.Mode().Perm()&(1<<(uint(7))) == 0 { //nolint:mnd // check write bit
		return false, errors.New("Write permission bit is not set for user. " + path)
	}

	return isWritableByOwner(path)
}

func Create(name string) (*os.File, error) {
	_, err := os.Stat(name)
	if err == nil {
		return nil, os.ErrExist
	}
	return os.Create(name)
}
