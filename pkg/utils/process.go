package utils

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"strconv"
)

// errCleanPID is the error returned when the
var errEmpty = errors.New("empty file")

// PID checks if other instances of the program is running (using the path to the .pid file provided).
// If so, it will exit with error code 0
// If not, it will create/update the .pid file provided
func PID(path string) error {
	// Read pid from path
	pid, err := readPID(path)
	if err != nil {
		// Create/overwrite pid file if it doesn't exist or is empty.
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, errEmpty) {
			return createPID(path)
		}

		// If other error, return
		return err
	}

	// Check if process with the pid provided exists
	exists, err := processExists(pid)
	if err != nil {
		return err
	}

	// Exit if exists, create pid file otherwise
	if exists {
		os.Exit(1)
	}
	return createPID(path)
}

// createPID creates/overwrites the path provided with a file that contains this program's PID
func createPID(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create pid file: %w", err)
	}
	defer f.Close()

	if _, err = f.WriteString(strconv.Itoa(os.Getpid()) + "\n"); err != nil {
		return fmt.Errorf("cannot write pid file: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("error closing pid file: %w", err)
	}
	return nil
}

// readPID reads the pid from the path provided
func readPID(path string) (int, error) {
	// Read data from path
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error reading pid file: %w", err)
	}

	// Check it's not empty
	if len(data) == 0 {
		return 0, errEmpty
	}

	// Remove last linebreak
	if data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	// Get the pid
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("error parsing PID: %w", err)
	}
	return pid, nil
}

// processExists checks if a process with the pid provided exists
func processExists(pid int) (bool, error) {
	if err := unix.Kill(pid, unix.Signal(0)); err != nil {
		if !errors.Is(err, unix.ESRCH) {
			return false, fmt.Errorf("error checking existence of process with pid %d: %w", pid, err)
		}
		return false, nil
	}
	return true, nil
}
