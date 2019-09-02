package utils

import (
	"fmt"
	"io"
	"os"
)

// ForEachInDirectory executes the function provided for each os.FileInfo found in the directory of the path provided.
// If the function returns an error, the execution will not stop and will be collected in an error slice.
func ForEachInDirectory(path string, fn func(fi os.FileInfo) error) []error {
	// Check if path exists, is readable, and is a directory
	fStat, err := os.Stat(path)
	if err != nil {
		return []error{err}
	}
	if !fStat.IsDir() {
		return []error{fmt.Errorf("%s is not a directory", path)}
	}

	// Open directory for reading
	f, err := os.Open(path)
	if err != nil {
		return []error{fmt.Errorf("error opening directory \"%s\": %s", path, err)}
	}
	defer f.Close()

	errList := make([]error, 0, 50)
	for {
		// List 1000 children each time.
		list, err := f.Readdir(1000)
		if err != nil {
			if err != io.EOF {
				return []error{fmt.Errorf("error listing directory \"%s\": %s", path, err)}
			}
			break
		}

		// Apply fn(os.FileInfo) for each child
		for _, fi := range list {
			if err = fn(fi); err != nil {
				errList = append(errList, err)
			}
		}
	}

	if len(errList) != 0 {
		return errList
	}
	return nil
}

