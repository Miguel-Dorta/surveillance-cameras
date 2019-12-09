package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"time"
)

const (
	user = "USER"
	pass = "PASS"
	path = "/auto.jpg"
)

type handler int

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	now := time.Now()

	receivedUsername, receivedPassword, ok := r.BasicAuth()
	if !ok || receivedUsername != user || receivedPassword != pass {
		fail("invalid request")
	}

	if r.Method != http.MethodGet {
		fail("invalid method: expected GET, found %s", r.Method)
	}

	if r.URL.Path != path {
		fail("invalid path: expected %s, found %s", path, r.URL.Path)
	}

	rw.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(rw, "%04d/%02d/%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func fail(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a...)
	os.Exit(1)
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	s := httptest.NewServer(new(handler))
	fmt.Println(s.URL)
	<-quit
	s.Close()
}
