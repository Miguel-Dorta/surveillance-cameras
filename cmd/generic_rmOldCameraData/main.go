package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	log *logolang.Logger
	path string
	oldestPreserve time.Time
)

func init() {
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelError

	var days int
	var verbose, version bool
	flag.StringVar(&path, "path", ".", "Path for removing old camera data")
	flag.IntVar(&days, "days", 30, "Preserve camera data more recent than ~ days")
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

	oldestPreserve = time.Now().AddDate(0, 0, days * -1)
}

func main() {
	if err := utils.IterateDir(path, func(f os.FileInfo) {
		if !f.IsDir() {
			return
		}
		iterateYears(filepath.Join(path, f.Name()))
	}); err != nil {
		log.Errorf("error listing path \"%s\": %s", path, err)
	}
}

func iterateYears(path string) {
	iterateDate(path, oldestPreserve.Year(), iterateMonths)
}

func iterateMonths(path string) {
	iterateDate(path, int(oldestPreserve.Month()), iterateDays)
}

func iterateDays(path string) {
	iterateDate(path, oldestPreserve.Day(), func(_ string) {})
}

func iterateDate(path string, compareTo int, doWithCoincidence func(path string)) {
	if err := utils.IterateDir(path, func(f os.FileInfo) {
		fPath := filepath.Join(path, f.Name())

		// Parse comparable
		comparable, err := strconv.Atoi(f.Name())
		if err != nil {
			log.Errorf("error parsing date from path \"%s\": %s", fPath, err)
			return
		}

		// Remove if it's an older comparable
		if comparable < compareTo {
			log.Infof("removing %s", fPath)
			if err := os.RemoveAll(fPath); err != nil {
				log.Errorf("error removing path \"%s\": %s", fPath, err)
			}
			return
		}

		// Do action if it's the same
		if comparable == compareTo {
			doWithCoincidence(fPath)
			return
		}

		// Skip if it's a future comparable

	}); err != nil {
		log.Errorf("error listing path \"%s\": %s", path, err)
	}
}
