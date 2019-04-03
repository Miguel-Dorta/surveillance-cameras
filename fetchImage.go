package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var camName, user, pass, url, path string
	if len(os.Args) == 6 {
		camName = os.Args[1]
		user = os.Args[2]
		pass = os.Args[3]
		url = os.Args[4]
		path = os.Args[5]
	} else if len(os.Args) == 4 {
		camName = os.Args[1]
		user = ""
		pass = ""
		url = os.Args[2]
		path = os.Args[3]
	} else {
		fmt.Println("Usage:    <camera-name> <user> <password> <url> <path-destiny>\n  [user] and [password] are optional")
		os.Exit(1)
	}

	seconds := time.Tick(time.Second)
	for {
		<- seconds
		date := time.Now()
		err := getImage(url, user, pass, camName, path, date)
		if err != nil {
			fmt.Printf("[%04d-%02d-%02d %02d:%02d:%02d] %s\n",
				date.Year(),
				date.Month(),
				date.Day(),
				date.Hour(),
				date.Minute(),
				date.Second(),
				err.Error(),
			)
		}
	}
}

func getImage(url, user, pass, camName, path string, date time.Time) error {
	// Request image
	client := &http.Client{ Timeout: time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}
	if user != "" {
		req.SetBasicAuth(user, pass)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Create and open file
	dirPath := filepath.Join(
		path,
		camName,
		fmt.Sprintf("%04d", date.Year()),
		fmt.Sprintf("%02d", date.Month()),
		fmt.Sprintf("%02d", date.Day()),
	)
	fPath := filepath.Join(dirPath, fmt.Sprintf("%02d-%02d-%02d%s", date.Hour(), date.Minute(), date.Second(), filepath.Ext(url)))
	f, err := os.Create(fPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error creating local file: %s", err.Error())
		}

		// If it fails because there is not container folder, creates it
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
		err = fmt.Errorf("error saving file: %s", err.Error())
	}
	return err
}
