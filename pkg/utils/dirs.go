package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// IterateDir iterates the dir provided applying the function provided to every file found.
func IterateDir(path string, fn func(f os.FileInfo)) error {
	// Check if path exists, is readable, and is a directory
	fStat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting info of \"%s\": %w", path, err)
	}
	if !fStat.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	// Open directory for reading
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening directory \"%s\": %w", path, err)
	}
	defer f.Close()

	for {
		// List 1000 children each time.
		list, err := f.Readdir(1000)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("error listing directory \"%s\": %w", path, err)
			}
			return nil
		}

		// Apply fn(os.FileInfo) for each child
		for _, fi := range list {
			fn(fi)
		}
	}

	return nil
}
