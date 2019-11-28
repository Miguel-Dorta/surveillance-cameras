package main

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type fileCheck struct {
	name string
	isDir, checked bool
}

func Test(t *testing.T) {
	fileList := randomFileList()
	tmpDirPath := "/tmp/listLargeDir_test"

	if err := os.Mkdir(tmpDirPath, 0777); err != nil {
		t.Skipf("cannot create tmp dir")
	}
	defer os.RemoveAll(tmpDirPath)

	// Create files
	for _, file := range fileList {
		path := filepath.Join(tmpDirPath, file.name)
		if file.isDir {
			if err := os.Mkdir(path, 0777); err != nil {
				t.Fatalf("cannot create dir \"%s\": %s", path, err)
			}
		} else {
			f, err := os.Create(path)
			if err != nil {
				t.Fatalf("cannot create dir \"%s\": %s", path, err)
			}
			f.Close()
		}
	}

	// This is a weird test >.<
	outStr := &strings.Builder{}
	errStr := &strings.Builder{}
	logOut = log.New(outStr, "", 0)
	logErr = log.New(errStr, "", 0)
	errs := ls(tmpDirPath)

	// There should be no errors
	if len(errs) != 0 {
		for _, err := range errs {
			t.Error(err)
		}
		t.FailNow()
	}
	if errStr.String() != "" {
		t.Fatalf("errors found in output:\n%s", errStr.String())
	}

	for _, outLine := range strings.Split(outStr.String(), "\n") {
		if len(outLine) < 26 {
			continue
		}
		name := outLine[:16]
		isDir, err := strconv.ParseBool(outLine[26:])
		if err != nil {
			t.Fatalf("malformed line: %s", outLine)
		}

		fileInMap, found := fileList[name]
		if !found {
			t.Fatalf("unknown file %s", name)
		}

		if fileInMap.isDir != isDir {
			t.Fatalf("wrong type in file %s", name)
		}

		fileInMap.checked = true
		fileList[name] = fileInMap
	}

	for _, file := range fileList {
		if !file.checked {
			t.Errorf("file %s not checked", file.name)
		}
	}
}

func randomFileList() map[string]fileCheck {
	numberOfFiles := 1234
	m := make(map[string]fileCheck, numberOfFiles)

	for i:=0; i<numberOfFiles; i++ {
		nameBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(nameBytes, uint64(rand.Int63()))
		name := hex.EncodeToString(nameBytes)

		m[name] = fileCheck{
			name:    name,
			isDir:   rand.Intn(2) == 1,
			checked: false,
		}
	}
	return m
}