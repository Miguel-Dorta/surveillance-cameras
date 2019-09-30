package main

import (
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/client"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	USAGE = "<url> <user> <password> <number-of-movements>"
	moveLeftUrl = "/cgi-bin/hi3510/ytleft.cgi"
	moveRightUrl = "/cgi-bin/hi3510/ytright.cgi"
)

var (
	url, user, password string

	logOut = log.New(os.Stdout, "", 0)
	logErr = log.New(os.Stderr, "[ERROR]", 0)
)

func init() {
	client.HttpClient = &http.Client{Timeout:time.Second}
}

func getArgs() (movements int) {
	internal.CheckSpecialArgs(os.Args, USAGE)
	if len(os.Args) != 5 {
		logOut.Fatalf("Usage:    %s %s  (use -h for help)\n", os.Args[0], USAGE)
	}

	movements, err := strconv.Atoi(os.Args[4])
	if err != nil {
		logErr.Fatalf("Invalid number \"%s\": %s\n", os.Args[4], err)
	}

	if movements < 1 {
		logErr.Fatalln("The number provided must be equal or higher than 1")
	}

	logErr.SetFlags(log.Ldate | log.Ltime)

	url = os.Args[1]
	user = os.Args[2]
	password = os.Args[3]
	return movements
}

func main() {
	movements := getArgs()
	for {
		for i:=0; i<movements; i++ {
			move(moveLeftUrl)
			time.Sleep(time.Second)
		}
		for i:=0; i<movements; i++ {
			move(moveRightUrl)
			time.Sleep(time.Second)
		}
	}
}

func move(subUrl string) {
	_, err := client.GetWithLogin(url + subUrl, user, password)
	if err != nil {
		logErr.Printf("error doing moving request: %s\n", err)
	}
}
