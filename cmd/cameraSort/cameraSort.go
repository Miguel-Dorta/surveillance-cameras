package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"io"
	"os"
	"path/filepath"
)

const USAGE = "<origin> <destination>"
var ErrIncorrectFormat = fmt.Errorf("incorrect format")

func main() {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
		os.Exit(1)
	}

	origin := os.Args[1]
	destiny := os.Args[2]

	err := internal.ForEachInDirectory(origin, func(fi os.FileInfo) error {
		fiName := fi.Name()
		if !fi.Mode().IsRegular() {
			return nil
		}

		camId, y, m, d, err := getFileNameParts(fiName)
		if err != nil {
			return nil
		}

		originPath := filepath.Join(origin, fiName)
		destinyPath := filepath.Join(destiny, camId, y, m, d, fiName)

		err = os.Rename(originPath, destinyPath)
		if err != nil {
			if _, ok := err.(*os.LinkError); !ok {
				return fmt.Errorf("Error moving %s to %s: %s\n", originPath, destinyPath, err.Error())
			}

			err = copyFileBetweenDisks(originPath, destinyPath)
			if err != nil {
				return fmt.Errorf("Error copying %s to %s: %s\n", originPath, destinyPath, err.Error())
			}

			err = os.Remove(originPath)
			if err != nil {
				return fmt.Errorf("Error removing %s from source: %s\n", originPath, err.Error())
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf(":: Errors found:\n%s\n", err.Error())
		os.Exit(1)
	}
}

// Gets CamID, Year, Month and Day from a string like CAMIDyyyyMMdd
func getFileNameParts(s string) (camId, y, m, d string, err error) {
	name := []byte(s)
	err = ErrIncorrectFormat

	var c byte
	var camIdSize int
	for camIdSize, c = range name {
		if c > 47 && c < 58 { // Is it a number?
			break
		}
		if c < 65 || c > 122 || (c > 90 && c < 97) { // Is it NOT a letter? [a-zA-Z]
			return
		}
	}

	// Check if there's CamID and a full date
	if camIdSize == 0 || len(name) < camIdSize + 8 {
		return
	}

	for i := camIdSize; i < camIdSize + 8; i++ {
		if name[i] < 48 || name[i] > 57 { // Is it NOT a number?
			return
		}
	}
	return s[:camIdSize], s[camIdSize:camIdSize+4], s[camIdSize+4:camIdSize+6], s[camIdSize+6:camIdSize+8], nil
}

func copyFileBetweenDisks(origin, destiny string) (err error) {
	originFile, err := os.Open(origin)
	if err != nil {
		return
	}
	defer originFile.Close()

	destinyFile, err := os.Create(destiny)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error creating destiny file: %s", err.Error())
		}

		// If it fails because there is not container folders, creates them
		err := os.MkdirAll(getPathToFile(destiny), 0700)
		if err != nil {
			return fmt.Errorf("error creating container folders: %s", err.Error())
		}

		destinyFile, err = os.Create(destiny)
		if err != nil {
			return fmt.Errorf("error creating destiny file: %s", err.Error())
		}
	}
	defer destinyFile.Close()

	if _, err = io.CopyBuffer(destinyFile, originFile, make([]byte, 100*1024)); err != nil {
		return
	}

	if err = destinyFile.Close(); err != nil {
		return
	}

	return nil
}

func getPathToFile(path string) string {
	for i := len(path) - 1; i > 0; i-- {
		if path[i] == os.PathSeparator {
			return path[:i]
		}
	}
	return path
}
