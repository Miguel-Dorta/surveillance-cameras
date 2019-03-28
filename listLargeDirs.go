package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage:    ./listLargeDirs [path-optional]")
		os.Exit(1)
	}

	var path string
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		path, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting working directory: %s\nTry using ./listLargeDirs <absolutePath>\n", err.Error())
			os.Exit(1)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening directory: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	errCounter := 0
	for {
		list, err := f.Readdir(1000)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			errCounter++
			fmt.Printf("Error listing files [try %d of 10]: %s\n", errCounter, err.Error())

			if errCounter >= 10 {
				break
			}
		}

		for _, fi := range list {
			fmt.Printf("%s @ IsDir? %t\n", fi.Name(),fi.IsDir())
		}
	}
	os.Exit(errCounter)
}
