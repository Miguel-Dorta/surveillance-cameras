package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:    ./rmOldCamera <path> <days-to-preserve>")
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

					// Open directory
					f, err := os.Open(camYearMonthDayPath)
					if err != nil {
						if os.IsNotExist(err) {
							continue
						}
						fmt.Printf("Error reading directory \"%s\": %s\nSkipping directory...\n", camYearMonthDayPath, err.Error())
					}

					errCounter := 0
					for {
						// Read next 1000 files
						fileList, err := f.Readdir(1000)
						if err != nil {
							if err.Error() == "EOF" {
								break
							}

							if errCounter < 10 {
								errCounter++
								fmt.Printf("[try %d of 10] Error reading directory \"%s\": %s\n", errCounter, camYearMonthDayPath, err.Error())
								continue
							} else {
								fmt.Println("[SEVERE] Error accumulation. Skipping directory...")
								errCounter = 0
								break
							}
						}

						// Remove files
						for _, fileToRm := range fileList {
							err = os.Remove(filepath.Join(camYearMonthDayPath, fileToRm.Name()))
							if err != nil {
								if os.IsNotExist(err) {
									continue
								}
								fmt.Printf("Error removing file \"%s\": %s\n", filepath.Join(camYearMonthDayPath, fileToRm.Name()), err.Error())
							}
						}
					}
					f.Close()
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
