package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const USAGE = "<camera-name> <user-optional> <password-optional> <url> <path-destiny>"

var camName, user, pass, url, path string

func main() {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) == 6 {
		camName = os.Args[1]
		user = os.Args[2]
		pass = os.Args[3]
		url = os.Args[4]
		path = os.Args[5]
	} else if len(os.Args) == 4 {
		camName = os.Args[1]
		url = os.Args[2]
		path = os.Args[3]
		user = ""
		pass = ""
	} else {
		fmt.Printf("Usage:    %s %s (use -h for help)\n", os.Args[0], USAGE)
		os.Exit(1)
	}

	var (
		dirPath string
		tomorrow int64
		now time.Time
		err error
		fileExt = filepath.Ext(url)
		seconds = time.Tick(time.Second)
	)
	for {
		<- seconds
		now = time.Now().UTC()

		// If date changes, update dirPath
		if now.Unix() >= tomorrow {
			tomorrow = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(time.Hour * 24).Unix()
			dirPath = filepath.Join(
				path,
				camName,
				fmt.Sprintf("%04d", now.Year()),
				fmt.Sprintf("%02d", now.Month()),
				fmt.Sprintf("%02d", now.Day()),
			)
		}

		err = getImage(dirPath, fileExt, now)
		if err != nil {
			fmt.Printf("[%04d-%02d-%02d %02d:%02d:%02d] %s\n",
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute(),
				now.Second(),
				err.Error(),
			)
		}
	}
}

func getImage(dirPath, fileExt string, now time.Time) error {
	// Request image
	client := &http.Client{ Timeout: time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}
	if user != "" || pass != "" {
		req.SetBasicAuth(user, pass)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Create and open file
	fPath := filepath.Join(dirPath, fmt.Sprintf("%02d-%02d-%02d%s", now.Hour(), now.Minute(), now.Second(), fileExt))
	f, err := os.Create(fPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error creating local file: %s", err.Error())
		}

		// If it fails because there is not container folders, creates them
		err := os.MkdirAll(dirPath, 0700)
		if err != nil {
			return fmt.Errorf("error creating container folder: %s", err.Error())
		}

		f, err = os.Create(fPath)
		if err != nil {
			return fmt.Errorf("error creating local file: %s", err.Error())
		}
	}
	defer f.Close()

	// Save data
	_, err = io.CopyBuffer(f, resp.Body, make([]byte, 100 * 1024))
	if err != nil {
		return fmt.Errorf("error saving file: %s", err.Error())
	}

	if err = f.Close(); err != nil {
		return err
	}
	return nil
}
