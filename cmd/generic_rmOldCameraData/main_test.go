package main_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

const (
	programName = "generic_rmOldCameraData"
)

type date struct {
	cam              string
	year, month, day int
}

func now() *date {
	y, m, d := time.Now().Date()
	return &date{
		year:  y,
		month: int(m),
		day:   d,
	}
}

func (d1 *date) isOlderThan(d2 *date) bool {
	d11 := time.Date(d1.year, time.Month(d1.month), d1.day, 0, 0, 0, 0, time.Local)
	d22 := time.Date(d1.year, time.Month(d1.month), d1.day, 0, 0, 0, 0, time.Local)
	return d11.Before(d22)
}

func Test(t *testing.T) {
	var (
		tmpDir       = filepath.Join(os.TempDir(), programName)
		dataDir      = filepath.Join(tmpDir, "data")
		buildPath    = filepath.Join(tmpDir, "build", programName)
		numberOfCams = 4
	)

	// Remove previous test folder
	t.Log("Removing previous test folder")
	rmRf(tmpDir, t)

	// Build binary
	t.Log("Building binary")
	goBuild(filepath.Join(getEnv("GOPATH", t), "src", "github.com", "Miguel-Dorta", "surveillance-cameras", "cmd", programName), buildPath, t)

	// Save list of preserved dates for checking
	preservedDates := make(map[date]bool, numberOfCams*200*12*31)
	now := now()

	// Make test data
	t.Log("Making test data...")
	for c := 0; c <= numberOfCams; c++ {
		pathCam := filepath.Join(dataDir, fmt.Sprintf("C%03d", c))

		for y := 1900; y <= 2100; y++ {
			pathCamYear := filepath.Join(pathCam, strconv.Itoa(y))

			for m := 1; m <= 12; m++ {
				pathCamYearMonth := filepath.Join(pathCamYear, fmt.Sprintf("%02d", m))

				for d := 1; d <= 31; d++ {
					pathCamYearMonthDay := filepath.Join(pathCamYearMonth, fmt.Sprintf("%02d", d))

					mkdirP(pathCamYearMonthDay, t)
					touch(filepath.Join(pathCamYearMonthDay, "testfile"), t)

					dat := &date{cam: fmt.Sprintf("C%03d", c), year: y, month: m, day: d}
					if now.isOlderThan(dat) {
						preservedDates[*dat] = false
					}
				}
			}
		}
	}

	// Execute it
	t.Log("Executing binary")
	cmd := exec.Command(buildPath, "-path", dataDir, "-days", "0")
	if output, err := cmd.CombinedOutput(); err != nil || len(output) != 0 {
		t.Fatalf("error found executing binary: %s\nOUTPUT: %s", err, string(output))
	}

	// Check result
	t.Log("Checking for correct result (1/2)")
	// Check that there's not files that should have been deleted
	if err := filepath.Walk(dataDir, func(path string, _ os.FileInfo, err error) error {
		// Check error
		if err != nil {
			t.Errorf("cannot check path %s: %s", path, err)
			return nil
		}

		// Split relative path to its parts
		parts := filepath.SplitList(path)
		if len(parts) != 4 {
			return nil
		}

		// Convert numeric parts into ints
		y, err := strconv.Atoi(parts[1])
		if err != nil {
			t.Errorf("error parsing year in path %s: %s", path, err)
			return nil
		}
		m, err := strconv.Atoi(parts[2])
		if err != nil {
			t.Errorf("error parsing month in path %s: %s", path, err)
			return nil
		}
		d, err := strconv.Atoi(parts[3])
		if err != nil {
			t.Errorf("error parsing day in path %s: %s", path, err)
			return nil
		}

		// Check if the date was saved in preservedDates, that is, if it should be preserved or not.
		dat := date{cam: parts[0], year: y, month: m, day: d}
		if _, ok := preservedDates[dat]; !ok {
			t.Errorf("path %s was not supposed to be preserved", path)
			return nil
		}
		preservedDates[dat] = true
		return nil
	}); err != nil {
		t.Errorf("error walking tree directory: %s", err)
	}

	// Check if files that should have been preserved were deleted.
	t.Log("Checking for correct result (2/2)")
	for k, v := range preservedDates {
		if !v {
			t.Errorf("path %s was not preserved!", filepath.Join(k.cam, strconv.Itoa(k.year), strconv.Itoa(k.month), strconv.Itoa(k.day)))
		}
	}

	// Clean up
	t.Log("Clean up")
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Errorf("cannot remove tmp dir %s: %s", tmpDir, err)
	}
}

func getEnv(key string, t *testing.T) string {
	value := os.Getenv(key)
	if value == "" {
		t.Fatalf("env variable %s not defined", key)
	}
	return value
}

func rmRf(path string, t *testing.T) {
	if err := os.RemoveAll(path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("error removing path %s: %s", path, err)
		}
	}
}

func goBuild(src, to string, t *testing.T) {
	cmd := exec.Command("go", "build", "-o", to, src)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("error building binary: %s\n%s", err, string(output))
	}
}

func mkdirP(path string, t *testing.T) {
	if err := os.MkdirAll(path, 0777); err != nil {
		t.Fatalf("cannot create dir %s: %s", path, err)
	}
}

func touch(path string, t *testing.T) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("cannot create file %s: %s", path, err)
	}

	if err = f.Close(); err != nil {
		t.Fatalf("cannot close file %s: %s", path, err)
	}
}
