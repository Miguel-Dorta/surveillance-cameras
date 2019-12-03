package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const numberOfFiles = 300

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s <program-stdout-file> <file-list-path>\n", os.Args[0])
		os.Exit(1)
	}

	actualFiles := parseStdout(os.Args[1])
	expectedFiles := listedFiles(os.Args[2])

	for name, isDir := range actualFiles {
		b, exists := expectedFiles[name]
		if !exists {
			fail("found unexpected file (%s)", name)
		}

		if isDir != b {
			fail("in file (%s), the isDir found (%t) is different of the one expected (%t)", name, isDir, b)
		}
	}

	for name, isDir := range expectedFiles {
		b, exists := actualFiles[name]
		if !exists {
			fail("not found expected file (%s)", name)
		}

		if isDir != b {
			fail("in file (%s), the isDir found (%t) is different of the one expected (%t)", name, isDir, b)
		}
	}
}

func parseStdout(path string) map[string]bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fail("cannot read output file in path %s: %s", path, err)
	}

	files := make(map[string]bool, numberOfFiles)

	for _, line := range bytes.Split(data, []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}

		lineStr := string(line)

		parts := strings.Split(lineStr, " @ IsDir? ")
		if len(parts) != 2 {
			fail("more or less matches than expected in string \"%s\"", lineStr)
		}

		isDir, err := strconv.ParseBool(parts[1])
		if err != nil {
			fail("error parsing bool in line (%s): %s", lineStr, err)
		}

		files[parts[0]] = isDir
	}
	return files
}

func listedFiles(path string) map[string]bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fail("cannot read file list from file %s: %s", path, err)
	}

	files := make(map[string]bool, numberOfFiles)

	previousSemicolon, previousComma := -1, -1
	for i, b := range data {
		if b == ',' {
			previousComma = i
			continue
		}

		if b != ';' {
			continue
		}

		isDir, err := strconv.ParseBool(string(data[previousComma + 1:i]))
		if err != nil {
			fail("error parsing boolean from listed files: %s", err)
		}

		files[string(data[previousSemicolon + 1:previousComma])] = isDir
		previousSemicolon = i
	}
	return files
}

func fail(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a...)
	os.Exit(1)
}
