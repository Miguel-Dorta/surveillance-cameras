package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/httpClient"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

var (
	camName, user, pass, url, path string
	destinyDir, fileExtension      string
	printVersion                   bool
	tomorrow                       time.Time
)

func init() {
	flag.StringVar(&camName, "camera-name", "", "Sets the camera name/ID")
	flag.StringVar(&user, "user", "", "Username for login [optional]")
	flag.StringVar(&pass, "password", "", "Password for login [optional]")
	flag.StringVar(&url, "url", "", "URL for fetching the images")
	flag.StringVar(&path, "path", "", "Path for saving the images")
	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.BoolVar(&printVersion, "V", false, "Print version and exit")

	httpClient.Client.Timeout = time.Second
}

func checkFlags() {
	flag.Parse()

	if printVersion {
		fmt.Println(internal.Version)
		os.Exit(1)
	}

	if camName == "" || url == "" || path == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Arguments \"--camera-name\", \"--url\" and \"--path\" are required")
		os.Exit(1)
	}

	fileExtension = filepath.Ext(url)
	if fileExtension == "" {
		_, _ = fmt.Fprintln(os.Stderr, "WARNING: URL does not contain extension. Files will be saved without it.")
	}
}

func main() {
	checkFlags()

	var (
		seconds = time.Tick(time.Second)
		quit    = make(chan os.Signal, 2)
	)
	signal.Notify(quit, unix.SIGTERM, unix.SIGINT)

MainLoop:
	for range seconds {
		select {
		case <-quit:
			break MainLoop
		default:
		}

		requestTime := time.Now()
		updateDestinyDir(requestTime)

		if err := httpClient.GetFileWithLogin(url, user, pass, getNewFilePath(path, requestTime)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error downloading image: %s", err)
		}
	}
}

func updateDestinyDir(requestTime time.Time) {
	// If request time is before tomorrow, do not update destinyDir
	if requestTime.Before(tomorrow) {
		return
	}

	// Set tomorrow as the first instant of the next day of requestTime
	tomorrow = time.Date(requestTime.Year(), requestTime.Month(), requestTime.Day(), 0, 0, 0, 0, requestTime.Location())
	destinyDir = filepath.Join(
		path,
		camName,
		fmt.Sprintf("%04d", requestTime.Year()),
		fmt.Sprintf("%02d", requestTime.Month()),
		fmt.Sprintf("%02d", requestTime.Second()),
	)

	// Get file info (for checking that it exists and it's a dir)
	stat, err := os.Stat(destinyDir)
	if err == nil {
		return // Done, directory exists
	}
	if !os.IsNotExist(err) {
		panic(err) // If unknown error, panic
	}

	// If it's not dir, panic
	if !stat.IsDir() {
		panic("destiny path " + destinyDir + " MUST be a directory")
	}

	// Create new destinyDir. If cannot create it, panic
	if err = os.MkdirAll(destinyDir, 0755); err != nil {
		panic(fmt.Sprintf("cannot create destiny dir in \"%s\": %s", destinyDir, err))
	}
}

func getNewFilePath(saveTo string, requestTime time.Time) string {
	return filepath.Join(saveTo, fmt.Sprintf(
		"%02d-%02d-%02d%s",
		requestTime.Hour(),
		requestTime.Minute(),
		requestTime.Second(),
		fileExtension,
	))
}
