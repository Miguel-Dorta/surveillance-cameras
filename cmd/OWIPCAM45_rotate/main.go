package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/client"
	"net/http"
	"os"
	"time"
)

const (
	moveLeftUrl = "/cgi-bin/hi3510/ytleft.cgi"
	moveRightUrl = "/cgi-bin/hi3510/ytright.cgi"
)

var (
	url, user, pass string
	numberOfMovements int
	log *logolang.Logger
)

func init() {
	client.HttpClient = &http.Client{Timeout:time.Second}
	log = logolang.NewLogger()
	log.Level = logolang.LevelError

	var verbose, version bool
	flag.StringVar(&url, "url", "", "USAGE")
	flag.StringVar(&user, "user", "", "USAGE")
	flag.StringVar(&pass, "pass", "", "USAGE")
	flag.IntVar(&numberOfMovements, "movements", 10, "USAGE")
	flag.BoolVar(&verbose, "verbose", false, "USAGE")
	flag.BoolVar(&verbose, "v", false, "USAGE")
	flag.BoolVar(&version, "version", false, "USAGE")
	flag.BoolVar(&version, "V", false, "USAGE")
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
	if user == "" {
		log.Critical("invalid user")
		os.Exit(1)
	}
	if pass == "" {
		log.Critical("invalid password")
		os.Exit(1)
	}
	if numberOfMovements < 1 {
		log.Critical("invalid number of movements")
		os.Exit(1)
	}
}

func main() {
	for {
		for i:=0; i<numberOfMovements; i++ {
			log.Infof("Moving left (%d/%d)", i+1, numberOfMovements)
			move(moveLeftUrl)
			time.Sleep(time.Second)
		}
		for i:=0; i<numberOfMovements; i++ {
			log.Infof("Moving right (%d/%d)", i+1, numberOfMovements)
			move(moveRightUrl)
			time.Sleep(time.Second)
		}
	}
}

func move(direction string) {
	_, err := client.GetWithLogin(url + direction, user, pass)
	if err != nil {
		log.Errorf("error doing move request: %s", err)
	}
}
