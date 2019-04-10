package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const USAGE = "<path> <days-to-preserve>"

func main() {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
		os.Exit(1)
	}

	path := os.Args[1]
	daysToPreserve, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error parsing number %s: %s\nIs it really a number?\n", os.Args[2], err.Error())
		os.Exit(1)
	}
	dateToRm := time.Now().AddDate(0, 0, daysToPreserve * -1)

	camList, err := listDir(path)
	if err != nil {
		fmt.Printf("Error reading directory \"%s\": %s\n", path, err.Error())
		os.Exit(1)
	}

	for _, cam := range camList {
		if !cam.IsDir() {
			continue
		}

		camPath := filepath.Join(path, cam.Name())
		yearToRm := dateToRm.Year() // Save in var to avoid unnecessary function calls
		for y := 1970; y <= yearToRm; y++ {
			camYearPath := filepath.Join(camPath, strconv.Itoa(y))

			var monthToRm int
			if y == yearToRm {
				monthToRm = int(dateToRm.Month())
			} else {
				monthToRm = 12
			}

			for m := 1; m <= monthToRm; m++ {
				camYearMonthPath := filepath.Join(camYearPath, strconv.Itoa(m))

				var dayToRm int
				if y == yearToRm && m == monthToRm {
					dayToRm = dateToRm.Day()
				} else {
					dayToRm = 32 //31+1
				}

				for d := 1; d < dayToRm; d++ {
					camYearMonthDayPath := filepath.Join(camYearMonthPath, strconv.Itoa(d))

					err = internal.ForEachInDirectory(camYearMonthDayPath, func(fi os.FileInfo) error {
						err = os.Remove(filepath.Join(camYearMonthDayPath, fi.Name()))
						if err != nil && !os.IsNotExist(err) {
							return fmt.Errorf("Error removing file \"%s\": %s\n", filepath.Join(camYearMonthDayPath, fi.Name()), err.Error())
						}
						return nil
					})
					if err != nil {
						if os.IsNotExist(err) {
							continue
						}
						fmt.Printf(":: Errors found removing content of \"%s\"\n%s\n:: Skipping directory...\n", camYearMonthDayPath, err.Error())
					}
				}
			}
		}
	}
}

func listDir(path string) (list []os.FileInfo, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	list, err = f.Readdir(-1)
	return
}
