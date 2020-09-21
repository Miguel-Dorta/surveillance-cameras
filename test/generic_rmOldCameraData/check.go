package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

var exitStatus = 0

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:    %s <path>\n", os.Args[0])
		os.Exit(1)
	}

	checkCamera(os.Args[1])
	os.Exit(exitStatus)
}

func checkCamera(path string) {
	expectedCameras := []string{"C000", "C001", "C002"}

	actualCameras, err := ListDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error listing path %s: %s\n", path, err)
		exitStatus = 1
		return
	}

	if !EqualsStringSlices(expectedCameras, actualCameras) {
		_, _ = fmt.Fprintf(os.Stderr, "expected camera list (%+v) is not what was found (%+v) in path %s\n", expectedCameras, actualCameras, path)
		exitStatus = 1
		return
	}

	for _, c := range actualCameras {
		checkYear(filepath.Join(path, c))
	}
}

func checkYear(path string) {
	now := time.Now()

	expectedYears := make([]string, 0, (2100-now.Year())+1)
	for y := now.Year(); y <= 2100; y++ {
		expectedYears = append(expectedYears, strconv.Itoa(y))
	}

	actualYears, err := ListDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error listing path %s: %s\n", path, err)
		exitStatus = 1
		return
	}

	if !EqualsStringSlices(expectedYears, actualYears) {
		_, _ = fmt.Fprintf(os.Stderr, "expected year list (%+v) is not what was found (%+v) in path %s\n", expectedYears, actualYears, path)
		exitStatus = 1
		return
	}

	for _, y := range actualYears {
		checkMonth(filepath.Join(path, y), y == strconv.Itoa(now.Year()), now)
	}
}

func checkMonth(path string, sameYear bool, now time.Time) {
	start := 1
	if sameYear {
		start = int(now.Month())
	}

	expectedMonths := make([]string, 0, 12)
	for i := start; i <= 12; i++ {
		expectedMonths = append(expectedMonths, fmt.Sprintf("%02d", i))
	}

	actualMonths, err := ListDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error listing path %s: %s\n", path, err)
		exitStatus = 1
		return
	}

	if !EqualsStringSlices(expectedMonths, actualMonths) {
		_, _ = fmt.Fprintf(os.Stderr, "expected month list (%+v) is not what was found (%+v) in path %s\n", expectedMonths, actualMonths, path)
		exitStatus = 1
		return
	}

	for _, m := range actualMonths {
		checkDay(filepath.Join(path, m), sameYear && m == fmt.Sprintf("%02d", now.Month()), now)
	}
}

func checkDay(path string, sameMonth bool, now time.Time) {
	start := 1
	if sameMonth {
		start = now.Day()
	}

	expectedDays := make([]string, 0, 31)
	for i := start; i <= 31; i++ {
		expectedDays = append(expectedDays, fmt.Sprintf("%02d", i))
	}

	actualDays, err := ListDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error listing path %s: %s\n", path, err)
		exitStatus = 1
		return
	}

	if !EqualsStringSlices(expectedDays, actualDays) {
		_, _ = fmt.Fprintf(os.Stderr, "expected day list (%+v) is not what was found (%+v) in path %s\n", expectedDays, actualDays, path)
		exitStatus = 1
		return
	}
}

func ListDir(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dirNames, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	sort.Strings(dirNames)
	return dirNames, err
}

func EqualsStringSlices(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
