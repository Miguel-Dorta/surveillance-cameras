package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
)

const USAGE = "[path-optional]"

func main() {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) > 2 {
		fmt.Printf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
		os.Exit(1)
	}

	var path string
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting working directory: %s\nTry using ./listLargeDirs <absolutePath>\n", err.Error())
			os.Exit(1)
		}
	}

	errs := utils.ForEachInDirectory(path, func(fi os.FileInfo) error {
		fmt.Printf("%s @ IsDir? %t\n", fi.Name(), fi.IsDir())
		return nil
	})
	if len(errs) != 0 {
		for _, err := range errs {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
