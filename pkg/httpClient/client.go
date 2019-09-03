package httpClient

import (
	"fmt"
	"net/http"
	"time"
)

var c = http.Client{Timeout: time.Second}

func GetLogin(url, user, pass string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	if user != "" || pass != "" {
		req.SetBasicAuth(user, pass)
	}
	return c.Do(req)
}
