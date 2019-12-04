package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileList struct {
	Files []*File
}

type File struct {
	Name string
	Date *Data
}

type Data struct {
	CamName string
	Year int
	Month int
	Day int
	Hour int
	Minute int
	Second int
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s <path> <file-list-path>\n", os.Args[0])
		os.Exit(1)
	}

	l := new(FileList)
	parse(os.Args[2], l)

	for _, f := range l.Files {
		data := readFileStr(filepath.Join(
			os.Args[1],
			f.Date.CamName,
			fmt.Sprintf("%04d", f.Date.Year),
			fmt.Sprintf("%02d", f.Date.Month),
			fmt.Sprintf("%02d", f.Date.Day),
			f.Name,
		))

		if data != f.Name {
			fail("info found (%s) is not what was expected (%s)", data, f.Name)
		}
	}
}

func readFileStr(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fail("readFileStr: %s", err)
	}
	return string(data)
}

func parse(path string, i interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fail("parse: error reading file %s: %s", path, err)
	}

	if err = json.Unmarshal(data, i); err != nil {
		fail("parse: error unmarshaling file %s: %s", path, err)
	}
}

func fail(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a...)
	os.Exit(1)
}
