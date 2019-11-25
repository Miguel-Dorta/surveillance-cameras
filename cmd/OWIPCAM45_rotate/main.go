package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/client"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
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
	log.Color = false
	log.Level = logolang.LevelError

	var (
		pidFile, camName string
		verbose, version bool
	)
	flag.StringVar(&url, "url", "", "URL of the camera")
	flag.StringVar(&user, "user", "", "User for login")
	flag.StringVar(&pass, "pass", "", "Password for login")
	flag.StringVar(&camName, "camera-name", "", "Sets the camera name/ID")
	flag.StringVar(&pidFile, "pid", "/run/OWIPCAM45_rotate_<camera-name>.pid", "Path for pid file")
	flag.IntVar(&numberOfMovements, "movements", 10, "Number of rotations")
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
	if camName == "" {
		log.Criticalf("invalid camera name")
		os.Exit(1)
	}

	// Check for other instances
	if pidFile == "/run/OWIPCAM45_rotate_<camera-name>.pid" {
		pidFile = "/run/OWIPCAM45_rotate_" + camName + ".pid"
	}
	if err := utils.PID(pidFile); err != nil {
		log.Criticalf("error checking for other instances: %s", err)
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
