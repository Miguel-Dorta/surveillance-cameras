package main

import (
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
)

var (
	path string
	log *logolang.Logger
)

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

	// Check special args
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-V":
			fallthrough
		case "--version":
			fmt.Println(internal.Version)
			os.Exit(0)

		case "-h":
			fallthrough
		case "--help":
			fmt.Printf(
				"Usage:    %s [path-optional]\n"+
					"  -h, --help       Show this help text.\n"+
					"  -V, --version    Display version and exits.\n",
				os.Args[0])
			os.Exit(0)
		}
	}

	if len(os.Args) > 2 {
		log.Criticalf("Usage:    %s [path-optional] (use -h for help)", os.Args[0])
		os.Exit(0)
	}

	path = "."
	if len(os.Args) == 2 {
		path = os.Args[1]
	}
}

//TODO rewrite tests
func main() {
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
