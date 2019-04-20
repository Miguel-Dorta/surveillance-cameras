package internal

import (
	"fmt"
	"os"
)

var Version string

func CheckSpecialArgs(args []string, usage string) {
	for _, arg := range args[1:] {
		if arg[0] != '-' {
			continue
		}

		if arg == "-V" || arg == "--version" {
			fmt.Println(Version)
		} else if arg == "-h" || arg == "--help" {
			fmt.Printf("Usage:    %s %s\n  -h, --help       Show this help text.\n  -V, --version    Display version and exits.\n", args[0], usage)
		} else {
			continue
		}
		os.Exit(0)
	}
}
