package main

import (
	"os"
	"time"
	"strconv"
	"fmt"
)

const (
	PS string = string(os.PathSeparator)
)

func main() {
	path := os.Args[1]
	preserveDays, err := strconv.Atoi(os.Args[2])
	checkError(err, 1)

	camList, err, debugNumber := getDirList(path, -1)
	checkError(err, debugNumber)

	dateToRm := time.Now().AddDate(0,0, preserveDays * -1)
	yRm, mM, dRm := dateToRm.Date()
	mRm := int(mM)

	for _, cam := range camList {
		if cam.IsDir() {
			camName := cam.Name()
			camPath := path + PS + camName

			for y := 1970; y <= yRm; y++ {
				yPath := camPath + PS + strconv.Itoa(y)

				for m := 1; m < 13; m++ {
					if y != yRm || m <= mRm {
						mPath := yPath + PS + fmt.Sprintf("%02d", m)

						for d := 1; d < 32; d++ {
							if m != mRm || d < dRm {
								dPath := mPath + PS + fmt.Sprintf("%02d", d)

								dFile, err := os.Open(dPath)
								checkError(err, 2)
								for {
									dlist, err := dFile.Readdir(1000)
									if err != nil {
										if err.Error() == "EOF" {
											break;
										}
										checkError(err, 3)
									}
									for _, snap := range dlist {
										err := os.Remove(dPath + PS + snap.Name())
										checkError(err, 4)
									}
								}
								dFile.Close()

								fmt.Printf("\rCamera: %s - Date: %d/%02d/%02d", camName, y, m, d)
							}
						}
					}
				}
			}
		}
	}
	print("\n")
}

func getDirList(path string, filesToReturn int) ([]os.FileInfo, error, int) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err, 101
	}
	defer f.Close()
	list, err := f.Readdir(filesToReturn)
	if err != nil {
		return nil, err, 102
	}
	return list, nil, 0
}

func checkError(err error, debugNumber int) {
	if err != nil {
		fmt.Printf("\nError %v\n%v\n", err, debugNumber)
		os.Exit(1)
	}
}
