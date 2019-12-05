package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s <path> <file-list-path>\n", os.Args[0])
		os.Exit(1)
	}

	l := newFileList(3000)

	// Make data dir
	mkdirP(os.Args[1])

	// Write files
	for _, f := range l.Files {
		echoTo(f.Name, filepath.Join(os.Args[1], f.Name))
	}

	// Serialize file list
	serialize(os.Args[2], l)
}

func newFileList(numberOfFiles int) *FileList {
	l := make([]*File, 0, numberOfFiles)
	for i:=0; i<numberOfFiles; i++ {
		l = append(l, randomFile())
	}
	return &FileList{Files:l}
}

func randomFile() *File {
	d := randomDate()
	name := fmt.Sprintf("%s%04d%02d%02d%02d%02d%02d%02d.jpg",
		d.CamName, d.Year, d.Month, d.Day, d.Hour, d.Minute, d.Second, rand.Intn(99) + 1,
	)
	return &File{
		Name: name,
		Date: d,
	}
}

func randomDate() *Data {
	return &Data{
		CamName: fmt.Sprintf("C%03d", rand.Intn(1000)),
		Year:   rand.Intn(151) + 1950,
		Month:  rand.Intn(12) + 1,
		Day:    rand.Intn(31) + 1,
		Hour:   rand.Intn(24),
		Minute: rand.Intn(60),
		Second: rand.Intn(60),
	}
}

func serialize(path string, i interface{}) {
	data, err := json.Marshal(i)
	if err != nil {
		fail("serialize: error marshaling interface: %s", err)
	}

	if err = ioutil.WriteFile(path, data, 0666); err != nil {
		fail("serialize: error writing file %s: %s", path, err)
	}
}

func mkdirP(path string) {
	if err := os.MkdirAll(os.Args[1], 0777); err != nil {
		fail("mkdirP: error creating dir %s: %s", path, err)
	}
}

func echoTo(msg, path string) {
	f, err := os.Create(path)
	if err != nil {
		fail("echoTo: error creating file %s: %s", path, err)
	}
	defer f.Close()

	if _, err := f.WriteString(msg); err != nil {
		fail("echoTo: error writing file %s: %s", path, err)
	}

	if err = f.Close(); err != nil {
		fail("echoTo: error closing file %s: %s", path, err)
	}
}

func fail(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a...)
	os.Exit(1)
}
