package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/cameras"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/client"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/html"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	foldersDir = "/sd/"
	videoDir   = "record000/"
)

var (
	url, user, pass, camName, destination string
	printVersion                          bool
	logErr                                = log.New(os.Stderr, "", 0)
)

func init() {
	flag.StringVar(&url, "url", "", "URL for fetching the videos")
	flag.StringVar(&user, "user", "", "User for login purposes")
	flag.StringVar(&pass, "password", "", "Password for login purposes")
	flag.StringVar(&destination, "path", "", "Path for saving the files")
	flag.StringVar(&camName, "camera-name", "", "Sets the camera name/ID")
	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.BoolVar(&printVersion, "V", false, "Print version and exit")

	client.HttpClient.Timeout = time.Hour
}

func parseFlags() {
	flag.Parse()

	if printVersion {
		fmt.Println(internal.Version)
		os.Exit(1)
	}

	if url == "" || user == "" || pass == "" || camName == "" || destination == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	parseFlags()

	linkVideos, err := getAllVideos(url, user, pass)
	if err != nil {
		logErr.Fatalf("cannot get a list of all videos: %s", err)
	}
	errFound := false
	for _, link := range linkVideos {
		parentPath, savingPath := getSavingPath(destination, camName, link.Text)

		// Make parent dirs if they don't exist
		if err = os.MkdirAll(parentPath, 0755); err != nil {
			logErr.Printf("cannot create parent directories of file \"%s\": %s\n", link.Text, err)
			errFound = true
			continue
		}

		// Check if file exists. If it does, skip it.
		if _, err := os.Stat(savingPath); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			logErr.Printf("error found getting info from file \"%s\": %s\n", savingPath, err)
			errFound = true
			continue
		}

		// Download video
		if err = client.GetFileWithLogin(url+link.HREF, user, pass, savingPath); err != nil {
			logErr.Printf("error saving file in path \"%s\": %s\n", savingPath, err)
			errFound = true
			continue
		}
	}
	if errFound {
		os.Exit(1)
	}
}

// getSavingPath returns where the file passed as argument (filename) must be saved (second string returned) as well as
// the path to its parent directory.
func getSavingPath(destination, camName, filename string) (parentPath, savingPath string) {
	y, m, d, rest := cameras.GetInfoFromFilenameOWIPCAM45(filename)
	parentPath = filepath.Join(destination, camName, "20"+y, m, d)
	return parentPath, filepath.Join(parentPath, rest)
}

// getAllVideos gets all the A elements that corresponds to videos from all the folders found in the page of the URL provided.
func getAllVideos(url, user, pass string) ([]html.A, error) {
	// Get page
	page, err := client.GetAllWithLogin(url+foldersDir, user, pass)
	if err != nil {
		return nil, fmt.Errorf("error getting page from URL \"%s\": %s", url, err)
	}

	// Get folders from page
	aList := html.GetAList(page)
	videoList := make([]html.A, 0, len(aList)*100)
	for _, a := range aList {
		if !cameras.IsValidFolderName(a.Text) {
			continue
		}
		videos, err := getVideoLinks(url+a.HREF+videoDir, user, pass)
		if err != nil {
			return nil, fmt.Errorf("error getting videos from URL \"%s\": %s", a.HREF, err)
		}
		videoList = append(videoList, videos...)
	}

	return videoList, nil
}

// getVideoLinks gets all the A elements of the page in the URL provided that contains a valid video name.
func getVideoLinks(url, user, pass string) ([]html.A, error) {
	page, err := client.GetAllWithLogin(url, user, pass)
	if err != nil {
		return nil, fmt.Errorf("error getting page from URL \"%s\": %s", url, err)
	}

	aList := html.GetAList(page)
	result := make([]html.A, 0, len(aList))
	for _, a := range aList {
		if cameras.IsValidVideoName(a.Text) {
			result = append(result, a)
		}
	}
	return result, nil
}
