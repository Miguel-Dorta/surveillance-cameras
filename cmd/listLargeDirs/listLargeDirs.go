package main

import (
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"log"
	"os"
)

const USAGE = "[path-optional]"

var (
	logOut = log.New(os.Stdout, "", 0)
	logErr = log.New(os.Stderr, "[ERROR]", 0)
)

func getArgs() (path string) {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) > 2 {
		logOut.Fatalf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
	}

	if len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			logErr.Fatalf("Cannot get working directory: %s\nTry using ./listLargeDirs <absolutePath>\n", err)
		}
	}
	return
}

func main() {
	path := getArgs()

	errs := ls(path)
	if len(errs) != 0 {
		for _, err := range errs {
			logErr.Println(err)
		}
		os.Exit(1)
	}
}

func ls(path string) []error {
	return utils.ForEachInDirectory(path, func(fi os.FileInfo) error {
		logOut.Printf("%s @ IsDir? %t\n", fi.Name(), fi.IsDir())
		return nil
	})
}
