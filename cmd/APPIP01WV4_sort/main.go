package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/si"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
	"path/filepath"
)

// nameParts represents the parts that of a filename that are useful for sorting
type nameParts struct {
	camName, y, m, d string
}

var (
	from, to string
	log *logolang.Logger
)

func init() {
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelError

	var verbose, version bool
	flag.StringVar(&from, "from", "", "Path to read the files")
	flag.StringVar(&to, "to", "", "Path to put the files")
	flag.StringVar(&si.Dir, "pid-directory", "/run", "Path to pid file's directory")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&version, "V", false, "Print version and exit")
	flag.Parse()

	if version {
		fmt.Println(internal.Version)
		os.Exit(0)
	}

	if verbose {
		log.Level = logolang.LevelInfo
	}

	if from == "" || to == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if si.Dir == "" {
		log.Critical("invalid pid directory")
		os.Exit(1)
	}

	if err := si.Register("APPIP01WV4_sort"); err != nil {
		if err == si.ErrOtherInstanceRunning {
			os.Exit(0)
		} else {
			log.Criticalf("error registering instance: %s", err)
			os.Exit(1)
		}
	}
}

func main() {
	errFound := false
	if err := utils.IterateDir(from, func(fi os.FileInfo) {
		fiPath := filepath.Join(from, fi.Name())

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
		destinyPath := filepath.Join(to, parts.camName, parts.y, parts.m, parts.d, fi.Name())
		log.Infof("copying file from %s to %s", fiPath, destinyPath)
		if err := utils.Move(fiPath, destinyPath); err != nil {
			log.Errorf("error copying file from %s to %s: %s", fiPath, destinyPath, err)
			errFound = true
			return
		}
	}); err != nil {
		log.Errorf("error iterating directory \"%s\": %s", from, err)
		errFound = true
	}

	if errFound {
		os.Exit(1)
	}
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
