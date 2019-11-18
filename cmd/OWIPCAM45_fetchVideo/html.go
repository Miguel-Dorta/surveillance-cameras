package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/client"
	"regexp"
)

// a represents a simple <a> tag in HTML containing an HREF and the text itself
type a struct {
	href, text string
}

// aRegex is the regex for parsing simple HTML <a> tags
var aRegex = regexp.MustCompile("<a href=\"([0-9A-Za-z_./]+)\">([0-9A-Za-z_./]+)</a>")

// getAllVideos gets all the A elements that corresponds to videos from all the folders found in the page of the URL provided.
func getAllVideos(url string) ([]a, error) {
	// Get page
	page, err := client.GetAllWithLogin(url+foldersDir, user, pass)
	if err != nil {
		return nil, fmt.Errorf("error getting page from URL \"%s\": %s", url, err)
	}

	// Get folders from page
	aList := getAList(page)
	videoList := make([]a, 0, len(aList)*100) // result video list
	// Iterate <a> links
	for _, a := range aList {
		// Get only the ones that corresponds to a folder
		if !isFolderName(a.text) {
			continue
		}

		videos, err := getVideoLinks(url+a.href+videoDir)
		if err != nil {
			return nil, fmt.Errorf("error getting videos from URL \"%s\": %s", a.href, err)
		}
		videoList = append(videoList, videos...)
	}

	return videoList, nil
}

// getVideoLinks gets all the <a> elements of the page in the URL provided that contains a valid video name.
func getVideoLinks(url string) ([]a, error) {
	page, err := client.GetAllWithLogin(url, user, pass)
	if err != nil {
		return nil, fmt.Errorf("error getting page from URL \"%s\": %s", url, err)
	}

	aList := getAList(page)
	result := make([]a, 0, len(aList))
	for _, a := range aList {
		if isValidVideoName(a.text) {
			result = append(result, a)
		}
	}
	return result, nil
}

// getAList gets a list of all the <a> tags in the data provided
func getAList(data []byte) []a {
	matches := aRegex.FindAllSubmatch(data, -1)
	result := make([]a, 0, len(matches))
	for _, match := range matches {
		if match[1] == nil || match[2] == nil {
			continue
		}
		result = append(result, a{
			href: string(match[1]),
			text: string(match[2]),
		})
	}
	return result
}
