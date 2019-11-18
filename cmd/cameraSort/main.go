package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
	"path/filepath"
)

var (
	from, to string
	log *logolang.Logger
)

func init() {
	log = logolang.NewLogger()
	log.Level = logolang.LevelError

	var verbose, version bool
	flag.StringVar(&from, "from", "", "Path to read the files")
	flag.StringVar(&to, "to", "", "Path to put the files")
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
}

// TODO rewrite tests
func main() {
	errFound := false
	if err := utils.IterateDir(from, func(f os.FileInfo) {
		fPath := filepath.Join(from, f.Name())

		// Omit if it's not a file
		if !f.Mode().IsRegular() {
			return
		}

		camID, year, month, day, err := getInfoFromFilename(f.Name())
		if err != nil {
			log.Errorf("cannot get info from file \"%s\": %s", fPath, err)
			errFound = true
			return
		}

		destination := filepath.Join(to, camID, year, month, day, f.Name())
		log.Infof("moving file from \"%s\" to \"%s\"", fPath, destination)
		if err = utils.Move(fPath, destination); err != nil {
			log.Errorf("error moving file from \"%s\" to \"%s\": %s", fPath, destination, err)
			errFound = true
			return
		}
	}); err != nil {
		log.Errorf("error listing path \"%s\": %s", from, err)
		errFound = true
	}

	if errFound {
		os.Exit(1)
	}
}

// Gets CamID, Year, Month and Day from a string like CMIDyyyyMMdd
func getInfoFromFilename(filename string) (camID, year, month, day string, err error) {
	if len(filename) < 12 {
		return "", "", "", "", errors.New("incorrect format: too short")
	}

	for _, b := range filename[4:12] {
		if b < '0' || b > '9' {
			return "", "", "", "", errors.New("incorrect format: cannot parse date")
		}
	}

	return filename[:4], filename[4:8], filename[8:10], filename[10:12], nil
}
