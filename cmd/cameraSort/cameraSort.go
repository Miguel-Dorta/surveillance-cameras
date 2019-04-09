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

	f, err := os.Open(origin)
	if err != nil {
		fmt.Printf("Error opening origin: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	errCounter := 0
	for {
		list, err := f.Readdir(1000)
		if err != nil {
			if err == io.EOF {
				break
			}

			if errCounter < 10 {
				fmt.Printf("[try %d of 10] Error reading origin: %s\n", errCounter, err.Error())
				errCounter++
				continue
			} else {
				fmt.Println("[SEVERE] Error accumulation")
				os.Exit(2)
			}
		}

		for _, e := range list {
			eName := e.Name()
			if !e.Mode().IsRegular() {
				continue
			}

			camId, y, m, d, err := getFileNameParts(eName)
			if err != nil {
				continue
			}

			originPath := filepath.Join(origin, eName)
			destinyPath := filepath.Join(destiny, camId, y, m, d, eName)

			err = os.Rename(originPath, destinyPath)
			if err != nil {
				if _, ok := err.(*os.LinkError); !ok {
					fmt.Printf("Error moving %s to %s: %s\n", originPath, destinyPath, err.Error())
					continue
				}

				err = copyFileBetweenDisks(originPath, destinyPath)
				if err != nil {
					fmt.Printf("Error copying %s to %s: %s\n", originPath, destinyPath, err.Error())
					continue
				}

				err = os.Remove(originPath)
				if err != nil {
					fmt.Printf("Error removing %s from source: %s\n", originPath, err.Error())
					continue
				}
			}
		}
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
		return
	}
	defer destinyFile.Close()

	_, err = io.CopyBuffer(destinyFile, originFile, make([]byte, 100*1024))
	return
}
