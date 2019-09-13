package main

import (
	"fmt"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/httpClient"
	"io/ioutil"
)

func main() {

}

func getPage(url, user, pass string) ([]byte, error) {
	resp, err := httpClient.GetLogin(url, user, pass)
	if err != nil {
		return nil, fmt.Errorf("error getting page from URL \"%s\": %s", url, err)
	}
	resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading page from URL \"%s\": %s", url, err)
	}

	return data, nil
}

