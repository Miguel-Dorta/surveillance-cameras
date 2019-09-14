package utils

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/httpClient"
	"golang.org/x/sys/unix"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var copyBuffer = make([]byte, 128 * 1024)

// Executes "mv" command creating the necessary directories in destiny
func Move(origin, destiny string) error {
	// Create parent directories if they don't exists
	if err := os.MkdirAll(filepath.Dir(destiny), 0755); err != nil {
		return fmt.Errorf("cannot create parent directories for path \"%s\": %s", destiny, err)
	}

	// Try to move the file by renaming it
	err := os.Rename(origin, destiny)
	if err == nil {
		return nil // Job done
	}

	// If it failed, and the error was not an invalid cross-device link, return original error
	if le, _ := err.(*os.LinkError); le.Err != unix.EXDEV {
		return err
	}

	// Copy file
	if err = Copy(origin, destiny); err != nil {
		return err
	}

	// Remove origin file after copying
	if err = os.Remove(origin); err != nil {
		return fmt.Errorf("error removing origin file in \"%s\": %s", origin, err)
	}

	return nil
}

// Copy will copy a file from the first path provided to the last path provided
func Copy(origin, destiny string) error {
	// Open origin file to read
	fOrigin, err := os.Open(origin)
	if err != nil {
		return fmt.Errorf("error opening origin file \"%s\": %s", origin, err)
	}
	defer fOrigin.Close()

	// Create and open destiny file to write
	fDestiny, err := os.Create(destiny)
	if err != nil {
		return fmt.Errorf("error creating destiny file \"%s\": %s", destiny, err)
	}
	defer fDestiny.Close()

	// Copy
	if _, err = io.CopyBuffer(fDestiny, fOrigin, copyBuffer); err != nil {
		return fmt.Errorf("error while copying file from \"%s\" to \"%s\": %s", origin, destiny, err)
	}

	// Close destiny checking for errors
	if err = fDestiny.Close(); err != nil {
		return fmt.Errorf("error closing destiny file in \"%s\": %s", destiny, err)
	}

	return nil
}

func GetFileWithLogin(url, user, pass, destination string, c http.Client) error {
	resp, err := httpClient.GetLogin(url, user, pass, c)
	if err != nil {
		return fmt.Errorf("error doing http request to URL \"%s\": %s", url, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("error creating file \"%s\": %s", destination, err)
	}
	defer f.Close()

	if _, err = io.CopyBuffer(f, resp.Body, copyBuffer); err != nil {
		return fmt.Errorf("error while saving file \"%s\": %s", destination, err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("error closing file \"%s\": %s", destination, err)
	}

	return nil
}
