package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// bufferSize is the size of the buffer for disk I/O.
// The buffer itself will be defined in each call to preserve the thread safety.
const bufferSize = 128 * 1024

// HttpClient is the client that this package will use
var HttpClient *http.Client

// GetWithLogin makes a HTTP GET request using basic auth
func GetWithLogin(url, user, pass string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	if user != "" || pass != "" {
		req.SetBasicAuth(user, pass)
	}
	return HttpClient.Do(req)
}

// GetFileWithLogin downloads a file from URL to destination using basic auth
func GetFileWithLogin(url, user, pass, destination string) error {
	resp, err := GetWithLogin(url, user, pass)
	if err != nil {
		return fmt.Errorf("error doing http request to URL \"%s\": %s", url, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("error creating file \"%s\": %s", destination, err)
	}
	defer f.Close()

	if _, err = io.CopyBuffer(f, resp.Body, make([]byte, bufferSize)); err != nil {
		// If failed, force close and try to remove it
		_ = f.Close()
		_ = os.Remove(destination)
		return fmt.Errorf("error while saving file \"%s\": %s", destination, err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("error closing file \"%s\": %s", destination, err)
	}

	return nil
}
