package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const numberOfFiles = 300

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s <path> <file-list-path>\n")
		os.Exit(1)
	}

	data := bytes.NewBuffer(make([]byte, 0, numberOfFiles * 10))
	files := generateRandomFileMap()
	mkdirP(os.Args[1])
	for name, isDir := range files {
		path := filepath.Join(os.Args[1], name)

		_, _ = fmt.Fprintf(data, "%s,%t;", name, isDir)

		if isDir {
			mkdirP(path)
		} else {
			touch(path)
		}
	}

	if err := ioutil.WriteFile(os.Args[2], data.Bytes(), 0666); err != nil {
		fail("error writing file list in path %s: %s", os.Args[2], err)
	}
}

func generateRandomFileMap() map[string]bool {
	m := make(map[string]bool, numberOfFiles)
	for i:=0; i<=numberOfFiles; i++ {
		name := randomString()
		if _, exists := m[name]; exists {
			i--
			continue
		}
		m[name] = rand.Uint32() % 2 == 1
	}
	return m
}

func randomString() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
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
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a...)
	os.Exit(1)
}
