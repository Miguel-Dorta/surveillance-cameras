package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	tmpDirPath := "/tmp/rmOldCamera_test"
	if err := os.Mkdir(tmpDirPath, 0777); err != nil {
		t.Fatalf("cannot create temporal dir %s", tmpDirPath)
	}
	// Not deleted for manual inspection

	createTestFiles(tmpDirPath, t)
	parseDaysToPreserve("30")
	errs := iterateCams(tmpDirPath)

	if len(errs) != 0 {
		for _, err := range errs {
			t.Error(err)
		}
		t.FailNow()
	}
}

func createTestFiles(path string, t *testing.T) {
	for c:=0; c<5; c++ {
		pathCam := filepath.Join(path, "C00" + strconv.Itoa(c))
		for y := 1970; y <= 2100; y++ {
			pathCamYear := filepath.Join(pathCam, strconv.Itoa(y))
			for m := 1; m <= 12; m++ {
				pathCamYearMonth := filepath.Join(pathCamYear, fmt.Sprintf("%02d", m))
				for d := 1; d <= 31; d++ {
					pathCamYearMonthDay := filepath.Join(pathCamYearMonth, fmt.Sprintf("%02d", d))
					if err := os.MkdirAll(pathCamYearMonthDay, 0777); err != nil {
						t.Errorf("cannot create day folder %s: %s", pathCamYearMonthDay, err)
						continue
					}
					f, err := os.Create(filepath.Join(pathCamYearMonthDay, "testFile"))
					if err != nil {
						t.Errorf("cannot create test file in day %s: %s", pathCamYearMonthDay, err)
						continue
					}
					f.Close()
				}
			}
		}
	}
}
