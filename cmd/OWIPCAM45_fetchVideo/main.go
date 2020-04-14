package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/si"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/client"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	url, user, pass, camName, destination string
	log *logolang.Logger
)

func init() {
	client.HttpClient = new(http.Client)
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelError

	var verbose, version bool
	flag.StringVar(&url, "url", "", "URL for fetching the videos")
	flag.StringVar(&user, "user", "", "User for login purposes")
	flag.StringVar(&pass, "password", "", "Password for login purposes")
	flag.StringVar(&destination, "path", ".", "Path for saving the files")
	flag.StringVar(&camName, "camera-name", "", "Sets the camera name/ID")
	flag.StringVar(&si.Dir, "pid-directory", "/run", "Path to pid file's directory")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&version, "V", false, "Print version and exit")
	flag.Parse()

	if version {
		fmt.Println(internal.Version)
		os.Exit(0)
	}

	if verbose {
		log.Level = logolang.LevelInfo
	}

	// Check for valid arguments
	if url == "" {
		log.Critical("invalid url")
		os.Exit(1)
	}
	if camName == "" {
		log.Critical("invalid camera name")
		os.Exit(1)
	}
	if si.Dir == "" {
		log.Critical("invalid pid dir")
		os.Exit(1)
	}

	// Check for other instances
	if err := si.Register("OWIPCAM45_fetchVideo_" + camName); err != nil {
		if err == si.ErrOtherInstanceRunning {
			os.Exit(0)
		} else {
			log.Criticalf("error registering program instance: %s", err)
			os.Exit(1)
		}
	}
}

func main() {
	for range time.NewTicker(time.Hour).C {
		fetchVideo()
	}
}

func fetchVideo() {
	client.HttpClient.Timeout = time.Second * 5

	linkVideos, err := getAllVideos(url)
	if err != nil {
		log.Criticalf("cannot get a list of all videos: %s", err)
		return
	}

	client.HttpClient.Timeout = time.Hour

	for _, link := range linkVideos {
		destination := getSavingPath(destination, camName, link.text)

		// Make parent dirs if they don't exist
		if err = os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
			log.Errorf("cannot create parent directories of file \"%s\": %s", link.text, err)
			continue
		}

		// Check if file exists. If it does, skip it.
		if _, err := os.Stat(destination); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			log.Errorf("cannot get information from file \"%s\": %s", destination, err)
			continue
		}

		// Download video
		log.Infof("downloading %s", link.text)
		if err = client.GetFileWithLogin(url+link.href, user, pass, destination); err != nil {
			log.Errorf("error downloading file \"%s\" in path \"%s\": %s", link.text, destination, err)
			continue
		}
	}
}

// getSavingPath returns where the file passed as argument (filename) must be saved (second string returned) as well as
// the path to its parent directory.
func getSavingPath(destination, camName, filename string) string {
	y, m, d, rest := getInfoFromFilename(filename)
	return filepath.Join(destination, camName, "20"+y, m, d, rest)
}
