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
	dateToRm := time.Now().UTC().AddDate(0, 0, daysToPreserve * -1)

	err = internal.ForEachInDirectory(path, func(cam os.FileInfo) error {
		if !cam.IsDir() {
			return nil
		}

		pathCam := filepath.Join(path, cam.Name())
		err = internal.ForEachInDirectory(pathCam, func(year os.FileInfo) error {
			y, err := strconv.Atoi(year.Name())
			if err != nil {
				return nil
			}
			pathCamYear := filepath.Join(pathCam, year.Name())

			if y < dateToRm.Year() {
				return os.RemoveAll(pathCamYear)
			}

			err = internal.ForEachInDirectory(pathCamYear, func(month os.FileInfo) error {
				m, err := strconv.Atoi(month.Name())
				if err != nil {
					return nil
				}
				pathCamYearMonth := filepath.Join(pathCamYear, month.Name())

				if m < int(dateToRm.Month()) {
					return os.RemoveAll(pathCamYearMonth)
				}

				err = internal.ForEachInDirectory(pathCamYearMonth, func(day os.FileInfo) error {
					d, err := strconv.Atoi(day.Name())
					if err != nil {
						return nil
					}
					pathCamYearMonthDay := filepath.Join(pathCamYearMonth, day.Name())

					if d < dateToRm.Day() {
						return os.RemoveAll(pathCamYearMonthDay)
					}
					return nil
				})
				return err
			})
			return err
		})
		return err
	})
}
