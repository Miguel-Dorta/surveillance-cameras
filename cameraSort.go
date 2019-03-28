package main

import (
	"os"
	"fmt"
	"regexp"
	"strings"
	"io/ioutil"
)

const (
	PS string = string(os.PathSeparator)
)

var processedFiles uint = 0

func main() {
	origin := os.Args[1]
	destiny := os.Args[2]
	re, err := regexp.Compile("([a-zA-z]+)(\\d{4})(\\d{2})(\\d{2}).*")

	f, err := os.Open(origin)
	checkError(err, 1)
	defer f.Close()

	for ;; {
		list, err := f.Readdir(1000)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("\nFiles processed: %v\n", processedFiles)
				break;
			}
			checkError(err, 2)
		}
		listLen := len(list)
		for index, fi := range list {
			fm := fi.Mode()
			name := fi.Name()
			if fm.IsRegular() && re.MatchString(name) {
				result := re.FindStringSubmatch(name)
				path := origin + PS + name
				suffix := result[1]
				y := result[2]
				m := result[3]
				d := result[4]
				finalPath := destiny + PS + suffix + PS + y + PS + m + PS + d + PS + name

				err = os.Rename(path, finalPath)
				if err != nil {
					if strings.HasSuffix(err.Error(), "invalid cross-device link") {
						file, err := ioutil.ReadFile(path)
						checkError(err, 4)
						err = ioutil.WriteFile(finalPath, file, os.FileMode(0777))
						checkError(err, 5)
						err = os.Remove(path)
						checkError(err, 6)
					} else {
						checkError(err, 3)
					}
				}
				processedFiles++
			}
			fmt.Printf("\r%v of %v", index + 1, listLen)
		}
		print("\n")
	}
}

func checkError(err error, debugNumber int) {
	if err != nil {
		fmt.Printf("\nFiles processed: %v\nError %v\n%v\n", processedFiles, debugNumber, err)
		os.Exit(1)
	}
}
