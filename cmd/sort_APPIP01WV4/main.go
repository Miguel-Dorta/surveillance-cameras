package main

import (
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
	"path/filepath"
)

// nameParts represents the parts that of a filename that are useful for sorting
type nameParts struct {
	camName, y, m, d string
}

const USAGE = "<origin> <destination>"

var log *logolang.Logger

func init() {
	log = logolang.NewLogger()
	log.Color = false
	log.Formatter = func(levelName, msg string) string {
		return fmt.Sprintf("%s: %s", levelName, msg)
	}
}

func main() {
	origin, destination := getArgs()

	errFound := false
	if err := utils.IterateDir(origin, func(fi os.FileInfo) {
		fiPath := filepath.Join(origin, fi.Name())

		// check if fi is a file
		if !fi.Mode().IsRegular() {
			log.Errorf("path %s is not a file", fiPath)
			errFound = true
			return
		}

		// get name parts
		parts := getNameParts(fi.Name())
		if parts == nil {
			log.Errorf("invalid filename %s", fi.Name())
			errFound = true
			return
		}

		// copy it to its destination
		destinyPath := filepath.Join(destination, parts.camName, parts.y, parts.m, parts.d, fi.Name())
		if err := utils.Move(fiPath, destinyPath); err != nil {
			log.Errorf("error copying file from %s to %s: %s", fiPath, destinyPath, err)
			errFound = true
			return
		}
	}); err != nil {
		log.Errorf("error iterating directory \"%s\": %s", origin, err)
		errFound = true
	}

	if errFound {
		os.Exit(1)
	}
}

func getArgs() (origin, destiny string){
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
		os.Exit(1)
	}
	return os.Args[1], os.Args[2]
}

// getNameParts gets the camera name, year, month and day from a filename like "MacAddress00(NAME)_0_YYYYMMDDhhmmss_number.jpg"
func getNameParts(s string) *nameParts {
	if len(s) < 41 {
		return nil
	}
	
	parts := nameParts{
		camName: s[13:17],
		y:       s[21:25],
		m:       s[25:27],
		d:       s[27:29],
	}

	if !isDecimalNumber(parts.y) || !isDecimalNumber(parts.m) || !isDecimalNumber(parts.d) {
		return nil
	}

	return &parts
}

// isDecimalNumber checks if the the string s correspond to an unsigned int in decimal
func isDecimalNumber(s string) bool {
	for _, b := range s {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
