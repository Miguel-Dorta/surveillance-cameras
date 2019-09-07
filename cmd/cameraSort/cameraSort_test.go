package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type file struct {
	camID, year, month, day string
	restOfFilename string
}

func (f *file) getFilename() string {
	return f.camID + f.year + f.month + f.day + f.restOfFilename
}

var validFiles = []file{
	{
		camID: "C000",
		year: "1970",
		month: "01",
		day: "01",
		restOfFilename: "000000.jpg",
	},
	{
		camID: "C001",
		year: "2019",
		month: "09",
		day: "07",
		restOfFilename: "211734.jpg",
	},
	{
		camID: "C999",
		year: "2099",
		month: "12",
		day: "31",
		restOfFilename: "235959.jpg",
	},
	{
		camID: "C042",
		year: "2014",
		month: "03",
		day: "21",
		restOfFilename: "182639.jpg",
	},
}

var invalidFiles = []string{
	"this-is-not-a-camfile.jpg",
	"OlderID20190907211734.jpg",
	"hi",
	"C003201904O2123456.png",
}

var folder = "C00120190101235959.jpg"

var expectedErrors = map[string]bool{
	"path /tmp/cameraSort_test/origin/C00120190101235959.jpg is not a file":true,
	"cannot parse file \"/tmp/cameraSort_test/origin/C003201904O2123456.png\": incorrect format: cannot parse date":true,
	"cannot parse file \"/tmp/cameraSort_test/origin/hi\": incorrect format: too short":true,
	"cannot parse file \"/tmp/cameraSort_test/origin/OlderID20190907211734.jpg\": incorrect format: cannot parse date":true,
	"cannot parse file \"/tmp/cameraSort_test/origin/this-is-not-a-camfile.jpg\": incorrect format: cannot parse date":true,
}

func createFiles(path string) error {
	// Create invalid files in path
	for _, invalidFile := range invalidFiles {
		f, err := os.Create(filepath.Join(path, invalidFile))
		if err != nil {
			return fmt.Errorf("error creating invalid file %s", invalidFile)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("cannot close invalid file %s", invalidFile)
		}
	}

	// Create valid files in path
	for _, validFile := range validFiles {
		f, err := os.Create(filepath.Join(path, validFile.getFilename()))
		if err != nil {
			return fmt.Errorf("error creating valid file %s", validFile.getFilename())
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("cannot close valid file %s", validFile.getFilename())
		}
	}

	if err := os.Mkdir(filepath.Join(path, folder), 0777); err != nil {
		return errors.New("cannot create folder" + folder)
	}

	return nil
}

func checkValidFiles(path string) error {
	for _, f := range validFiles {
		fPath := filepath.Join(path, f.camID, f.year, f.month, f.day, f.getFilename())
		_, err := os.Stat(fPath)
		if err != nil {
			return fmt.Errorf("error getting expected info from valid file %s sorted in %s", f.getFilename(), fPath)
		}
	}
	return nil
}

func checkInvalidFiles(path string) error {
	// Check for invalid files
	for _, f := range invalidFiles {
		_, err := os.Stat(filepath.Join(path, f))
		if err != nil {
			return fmt.Errorf("cannot find invalid file %s in original dir", f)
		}
	}
	if _, err := os.Stat(filepath.Join(path, folder)); err != nil {
		return fmt.Errorf("invalid folder %s not found", folder)
	}

	// Check that no valid file is left behind
	for _, f := range validFiles {
		if _, err := os.Stat(filepath.Join(path, f.getFilename())); err == nil || !os.IsNotExist(err) {
			return fmt.Errorf("valid file %s still exists in origin dir", f.getFilename())
		}
	}

	return nil
}

func Test(t *testing.T) {
	tmpDir := "/tmp/cameraSort_test"
	originDir := filepath.Join(tmpDir, "origin")
	resultDir := filepath.Join(tmpDir, "result")

	// Create tmp dirs
	defer os.RemoveAll(tmpDir)
	if err := os.MkdirAll(originDir, 0777); err != nil {
		t.Fatal("cannot create origin tmp dir")
	}
	if err := os.MkdirAll(resultDir, 0777); err != nil {
		t.Fatal("cannot create result tmp dir")
	}

	if err := createFiles(originDir); err != nil {
		t.Fatalf("error creating files for testing: %s", err)
	}

	errs := sortFiles(originDir, resultDir)
	if len(errs) != 0 {
		for _, err := range errs {
			if _, exists := expectedErrors[err.Error()]; exists {
				continue
			}
			t.Fatalf("unexpected error: %s", err)
		}
	}

	// Check if done correctly
	if err := checkValidFiles(resultDir); err != nil {
		t.Fatalf("unexpected result checking valid files: %s", err)
	}
	if err := checkInvalidFiles(originDir); err != nil {
		t.Fatalf("unexpected result checking invalid files: %s", err)
	}
}
