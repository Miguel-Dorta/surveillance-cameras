package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/cameras"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/html"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/httpClient"
	"io/ioutil"
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

	client.Client.Timeout = time.Hour
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
		pathToSave, err := createAndGetSavingPath(link.Text, camName, destination)
		if err != nil {
			logErr.Fatalf("cannot create parent directories of file \"%s\": %s", link.Text, err)
		}
		if err = client.GetFileWithLogin(url+link.HREF, user, pass, pathToSave); err != nil {
			logErr.Printf("error saving file in path \"%s\": %s", pathToSave, err)
			errFound = true
			continue
		}
	}
	if errFound {
		os.Exit(1)
	}
}

func createAndGetSavingPath(filename, camName, destination string) (string, error) {
	y, m, d, rest := cameras.GetInfoFromFilenameOWIPCAN45(filename)
	parentDirPath := filepath.Join(destination, camName, "20"+y, m, d)
	if err := os.MkdirAll(parentDirPath, 0755); err != nil {
		return "", fmt.Errorf("error creating parent directories: %s", err)
	}
	return filepath.Join(parentDirPath, rest), nil
}

func getAllVideos(url, user, pass string) ([]html.A, error) {
	// Get page
	page, err := getPage(url+foldersDir, user, pass)
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

func getVideoLinks(url, user, pass string) ([]html.A, error) {
	page, err := getPage(url, user, pass)
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

func getPage(url, user, pass string) ([]byte, error) {
	resp, err := client.GetWithLogin(url, user, pass)
	if err != nil {
		return nil, fmt.Errorf("error getting page from URL \"%s\": %s", url, err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading page from URL \"%s\": %s", url, err)
	}

	return data, nil
}
