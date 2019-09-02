package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/cameras"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"os"
	"path/filepath"
)

const USAGE = "<origin> <destination>"

func main() {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) != 3 {
		fmt.Printf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
		os.Exit(1)
	}

	origin := os.Args[1]
	destiny := os.Args[2]

	errs := utils.ForEachInDirectory(origin, func(fi os.FileInfo) error {
		fiPath := filepath.Join(origin, fi.Name())

		// Omit if it's not regular
		if !fi.Mode().IsRegular() {
			return fmt.Errorf("path %s is not a file", fiPath)
		}

		camID, year, month, day, err := cameras.GetInfoFromFilename(fi.Name())
		if err != nil {
			return fmt.Errorf("cannot parse file \"%s\": %s", fiPath, err)
		}

		destinyPath := filepath.Join(destiny, camID, year, month, day, fi.Name())
		if err = utils.Move(fiPath, destinyPath); err != nil {
			return fmt.Errorf("error moving file from \"%s\" to \"%s\": %s", fiPath, destinyPath, err)
		}

		return nil
	})

	if len(errs) != 0 {
		for _, err := range errs {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
