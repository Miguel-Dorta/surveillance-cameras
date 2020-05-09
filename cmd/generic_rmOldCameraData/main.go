package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

var (
	log        *logolang.Logger
	path       string
	days       int
	exitStatus = 0
)

func init() {
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelError

	var verbose, version bool
	flag.StringVar(&path, "path", ".", "Path for removing old camera data")
	flag.IntVar(&days, "days", 30, "Preserve ~ days of data")
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
		log.Level = logolang.LevelDebug
	}
}

func main() {
	camList, err := ioutil.ReadDir(path)
	if err != nil {
		log.Criticalf("error listing path \"%s\": %s", path, err)
		os.Exit(1)
	}

	for _, cam := range camList {
		camPath := filepath.Join(path, cam.Name())
		if !cam.IsDir() {
			log.Debugf("skipping file %s: not a directory", camPath)
			continue
		}
		iterateCam(camPath)
	}
	os.Exit(exitStatus)
}

func iterateCam(path string) {
	years, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("error iterating camera %s: %s", path, err)
		exitStatus = 1
		return
	}

	var dateList []int
	for _, year := range years {
		yearPath := filepath.Join(path, year.Name())
		if !year.IsDir() {
			log.Debugf("skipping file %s: not a directory", yearPath)
			continue
		}

		yearInt, err := strconv.Atoi(year.Name())
		if err != nil {
			log.Debugf("skipping file %s due to error parsing name into number: %s", yearPath, err)
			continue
		}
		dateList = append(dateList, iterateYear(yearPath, yearInt)...)
	}

	log.Debug("sorting dates")
	sort.Sort(sort.Reverse(sort.IntSlice(dateList)))

	var today = today()
	for i, date := range dateList {
		if date < today {
			log.Debugf("skipping %d days more recent than today", i)
			dateList = dateList[i:]
			break
		}
	}

	for i := days; i < len(dateList); i++ {
		y, m, d := getDate(dateList[i])
		dir := filepath.Join(path, strconv.Itoa(y), fmt.Sprintf("%02d", m), fmt.Sprintf("%02d", d))
		log.Debugf("removing %s", dir)
		if err := os.RemoveAll(dir); err != nil {
			log.Errorf("error removing dir \"%s\": %s", dir, err)
			exitStatus = 1
			continue
		}
	}

	log.Debug("removing empty directories")
	rmEmptyDirDeep(path, 2)
}

func iterateYear(path string, date int) []int {
	date <<= 4

	months, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("error listing dir \"%s\": %s", path, err)
		exitStatus = 1
		return nil
	}

	var list []int
	for _, month := range months {
		monthPath := filepath.Join(path, month.Name())
		if !month.IsDir() {
			log.Debugf("skipping file %s: not a directory", monthPath)
			continue
		}

		monthInt, err := strconv.Atoi(month.Name())
		if err != nil || monthInt < 1 || monthInt > 12 {
			log.Debugf("skipping file %s due to error parsing name into number or invalid month: %s", monthPath, err)
			continue
		}
		list = append(list, iterateMonth(monthPath, date|monthInt)...)
	}
	return list
}

func iterateMonth(path string, date int) []int {
	date <<= 5

	days, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("error listing dir \"%s\": %s", path, err)
		exitStatus = 1
		return nil
	}

	list := make([]int, 0, len(days))
	for _, day := range days {
		dayPath := filepath.Join(path, day.Name())
		if !day.IsDir() {
			log.Debugf("skipping file %s: not a directory", dayPath)
			continue
		}

		dayInt, err := strconv.Atoi(day.Name())
		if err != nil || dayInt < 1 || dayInt > 31 {
			log.Debugf("skipping file %s due to error parsing name into number or invalid day: %s", dayPath, err)
			continue
		}
		list = append(list, date|dayInt)
	}
	return list
}

func rmEmptyDirDeep(path string, depth int) {
	if depth < 1 {
		log.Debugf("trying to remove path %s", path)
		if err := rmDirIfEmpty(path); err != nil {
			log.Errorf("error removing dir \"%s\": %s", path, err)
			exitStatus = 1
			return
		}
		return
	}

	list, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("error listing dir \"%s\": %s", path, err)
		exitStatus = 1
		return
	}
	for _, child := range list {
		childPath := filepath.Join(path, child.Name())
		if !child.IsDir() {
			log.Debugf("skipping removing file %s: not a directory", childPath)
			continue
		}
		rmEmptyDirDeep(childPath, depth-1)
	}

	list, err = ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("error listing dir \"%s\": %s", path, err)
		exitStatus = 1
		return
	}
	if len(list) == 0 {
		if err := rmDirIfEmpty(path); err != nil {
			log.Errorf("error removing dir \"%s\": %s", path, err)
			exitStatus = 1
			return
		}
		return
	}
}

func rmDirIfEmpty(path string) error {
	if err := os.Remove(path); err != nil {
		if errors.Is(err, unix.ENOTEMPTY) || errors.Is(err, unix.EEXIST) {
			return nil
		}
		return err
	}
	return nil
}

func today() int {
	now := time.Now()
	return (now.Year() << 9) | (int(now.Month()) << 5) | now.Day()
}

func getDate(i int) (y, m, d int) {
	y = i >> 9
	m = (i >> 5) & 0b1111
	d = i & 0b11111
	return
}
