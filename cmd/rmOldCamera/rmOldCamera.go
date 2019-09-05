package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const USAGE = "<path> <days-to-preserve>"

var earliestTimeToPreserve time.Time

func parseDaysToPreserve(dtp string) {
	daysToPreserve, err := strconv.Atoi(dtp)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing number %s: %s\nIs it really a number?\n", dtp, err)
		os.Exit(1)
	}
	earliestTimeToPreserve = time.Now().AddDate(0, 0, daysToPreserve * -1)
}

func main() {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
		os.Exit(1)
	}

	path := os.Args[1]
	parseDaysToPreserve(os.Args[2])

	errs := iterateCams(path)
	if len(errs) != 0 {
		for _, err := range errs {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func iterateCams(path string) []error {
	var errs []error
	errs2 := utils.ForEachInDirectory(path, func(camDir os.FileInfo) error {
		// Omit if is not a directory
		if !camDir.IsDir() {
			return nil
		}

		errs = append(errs, iterateYears(filepath.Join(path, camDir.Name()))...)
		return nil
	})

	return append(errs, errs2...)
}

func iterateYears(path string) []error {
	var errs []error
	errs2 := utils.ForEachInDirectory(path, func(yearDir os.FileInfo) error {
		yearInt, err := strconv.Atoi(yearDir.Name())
		if err != nil {
			return fmt.Errorf("error parsing year in \"%s\": %s", path, err)
		}

		if yearInt < earliestTimeToPreserve.Year() {
			return os.RemoveAll(filepath.Join(path, yearDir.Name()))
		} else if yearInt > earliestTimeToPreserve.Year() {
			return nil
		}

		errs = append(errs, iterateMonths(filepath.Join(path, yearDir.Name()))...)
		return nil
	})

	return append(errs, errs2...)
}

func iterateMonths(path string) []error {
	var errs []error
	errs2 := utils.ForEachInDirectory(path, func(monthDir os.FileInfo) error {
		monthInt, err := strconv.Atoi(monthDir.Name())
		if err != nil {
			return fmt.Errorf("error parsing month in \"%s\": %s", path, err)
		}

		if monthInt < int(earliestTimeToPreserve.Month()) {
			return os.RemoveAll(filepath.Join(path, monthDir.Name()))
		} else if monthInt > int(earliestTimeToPreserve.Month()) {
			return nil
		}

		errs = append(errs, iterateDays(filepath.Join(path, monthDir.Name()))...)
		return nil
	})

	return append(errs, errs2...)
}

func iterateDays(path string) []error {
	return utils.ForEachInDirectory(path, func(dayDir os.FileInfo) error {
		dayInt, err := strconv.Atoi(dayDir.Name())
		if err != nil {
			return fmt.Errorf("error parsing day in \"%s\": %s", path, err)
		}

		if dayInt < earliestTimeToPreserve.Day() {
			return os.RemoveAll(filepath.Join(path, dayDir.Name()))
		}

		return nil
	})
}
