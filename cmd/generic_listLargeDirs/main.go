package main

import (
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
)

const USAGE = "[path-optional]"

var log *logolang.Logger

func init() {
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelInfo
	log.Formatter = func(levelName, msg string) string {
		if levelName != "ERROR" {
			return msg
		}
		return fmt.Sprintf("[%s] %s", levelName, msg)
	}
}

func getArgs() (path string) {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) > 2 {
		log.Criticalf("Usage:    %s %s (use -h for help)", os.Args[0], USAGE)
	}

	path = "."
	if len(os.Args) == 2 {
		path = os.Args[1]
	}
	return path
}

//TODO rewrite tests
func main() {
	path := getArgs()

	errFound := false
	if err := utils.IterateDir(path, func(f os.FileInfo) {
		log.Infof("%s @ IsDir? %t", f.Name(), f.IsDir())
	}); err != nil {
		log.Errorf("error listing path: %s", err)
		errFound = true
	}

	if errFound {
		os.Exit(1)
	}
}
