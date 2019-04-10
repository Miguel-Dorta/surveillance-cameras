package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func ForEachInDirectory(path string, fn func(fi os.FileInfo) error) error {
	fStat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fStat.IsDir() {
		return errors.New("error: %s is not a directory")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening directory: %s\n", err.Error())
	}
	defer f.Close()

	errList := make([]string, 0, 50)
	listErrors := 1
	for {
		list, err := f.Readdir(1000)
		if err != nil {
			if err == io.EOF {
				break
			}

			if listErrors < 10 {
				errList = append(errList, fmt.Sprintf("[try %d of 10] error listing directory: %s\n", listErrors, err.Error()))
				listErrors++
				continue
			} else {
				errList = append(errList, "[SEVERE] list errors accumulation")
				break
			}
		}

		for _, fi := range list {
			if err = fn(fi); err != nil {
				errList = append(errList, err.Error())
			}
		}
	}

	if len(errList) != 0 {
		var str strings.Builder
		str.Grow(1000)
		for _, errTxt := range errList {
			str.WriteString(errTxt)
		}
		err = errors.New(str.String())
	}
	return err
}
