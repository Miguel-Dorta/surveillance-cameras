package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:    %s <path>\n", os.Args[0])
		os.Exit(1)
	}

	for c := 0; c <= 5; c++ {
		cPath := filepath.Join(os.Args[1], fmt.Sprintf("C%03d", c))

		for y := 1950; y <= 2100; y++ {
			yPath := filepath.Join(cPath, strconv.Itoa(y))

			for m := 1; m <= 12; m++ {
				mPath := filepath.Join(yPath, fmt.Sprintf("%02d", m))

				for d := 1; d <= 31; d++ {
					dPath := filepath.Join(mPath, fmt.Sprintf("%02d", d))
					mkdirP(dPath)
					touch(filepath.Join(dPath, "testfile"))
				}
			}
		}
	}
}

func mkdirP(path string) {
	if err := os.MkdirAll(path, 0777); err != nil {
		fail("error creating directory %s: %s", path, err)
	}
}

func touch(path string) {
	f, err := os.Create(path)
	if err != nil {
		fail("error creating file %s: %s", path, err)
	}

	if err = f.Close(); err != nil {
		fail("error closing file %s: %s", path, err)
	}
}

func fail(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a)
	os.Exit(1)
}
