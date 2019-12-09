package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:    %s <path>\n", os.Args[0])
		os.Exit(1)
	}

	filesFound := false
	iterate(os.Args[1], func(camName os.FileInfo) {
		camPath := filepath.Join(os.Args[1], camName.Name())
		iterate(camPath, func(y os.FileInfo) {
			yPath := filepath.Join(camPath, y.Name())
			iterate(yPath, func(m os.FileInfo) {
				mPath := filepath.Join(yPath, m.Name())
				iterate(mPath, func(d os.FileInfo) {
					dPath := filepath.Join(mPath, d.Name())
					iterate(dPath, func(f os.FileInfo) {
						checkFile(
							filepath.Join(dPath, f.Name()),
							atoi(y.Name(), "year"),
							atoi(m.Name(), "month"),
							atoi(d.Name(), "day"),
							atoi(f.Name()[:2], "hour"),
							atoi(f.Name()[3:5], "minute"),
							atoi(f.Name()[6:8], "second"),
						)
						filesFound = true
					})
				})
			})
		})
	})

	if !filesFound {
		fail("no files found")
	}
}

func checkFile(path string, Y, M, D, h, m, s int) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fail("error reading file %s: %s", path, err)
	}

	if string(data) != fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", Y, M, D, h, m, s) {
		fail("info found (%s) wasn't expected in file %s", string(data), path)
	}
}

func atoi(number, varName string) int {
	i, err := strconv.Atoi(number)
	if err != nil {
		fail("error parsing %s (%s): %s", varName, number, err)
	}
	return i
}

func iterate(path string, f func(os.FileInfo)) {
	if err := utils.IterateDir(path, func(fi os.FileInfo) {
		f(fi)
	}); err != nil {
		fail("error iterating path %s", path)
	}
}

func fail(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
