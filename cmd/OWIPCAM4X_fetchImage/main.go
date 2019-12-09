package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/client"
	"github.com/Miguel-Dorta/surveillance-cameras/pkg/utils"
	"golang.org/x/sys/unix"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"
)

var (
	url, user, pass, camName, path, fileExtension string
	log *logolang.Logger
)

func init() {
	client.HttpClient = &http.Client{Timeout: time.Second}
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelError

	var (
		pidFile string
		verbose, version bool
	)
	flag.StringVar(&url, "url", "", "URL for fetching images")
	flag.StringVar(&user, "user", "", "Username for login")
	flag.StringVar(&pass, "password", "", "Password for login")
	flag.StringVar(&camName, "camera-name", "", "Sets the camera name/ID")
	flag.StringVar(&path, "path", "", "Path to save images fetched")
	flag.StringVar(&pidFile, "pid", "/run/OWIPCAM4X_fetchImage_<camera-name>.pid", "Path to pid file")
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

	// Check correct args
	switch "" {
	case url:
		log.Criticalf("invalid url")
		os.Exit(1)
	case user:
		log.Criticalf("invalid username")
		os.Exit(1)
	case pass:
		log.Criticalf("invalid password")
		os.Exit(1)
	case camName:
		log.Criticalf("invalid camera-name")
		os.Exit(1)
	case path:
		log.Criticalf("invalid path destination")
		os.Exit(1)
	case pidFile:
		log.Criticalf("invalid pid path")
		os.Exit(1)
	}
	fileExtension = filepath.Ext(url)

	// Check for other instances running
	if pidFile == "/run/OWIPCAM4X_fetchImage_<camera-name>.pid" {
		pidFile = "/run/OWIPCAM4X_fetchImage_" + camName + ".pid"
	}
	if err := utils.PID(pidFile); err != nil {
		log.Criticalf("error checking for other instances: %s", err)
		os.Exit(1)
	}
}

func main() {
	var (
		seconds = time.NewTicker(time.Second).C
		quit = make(chan os.Signal, 2)
	)
	signal.Notify(quit, unix.SIGTERM, unix.SIGINT)

	for {
		select {
		case <-seconds:
			fetchImage()
		case <-quit:
			return
		}
	}
}

func fetchImage() {
	now := time.Now()
	destination := filepath.Join(
		path,
		camName,
		strconv.Itoa(now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()),
		fmt.Sprintf("%02d-%02d-%02d%s", now.Hour(), now.Minute(), now.Second(), fileExtension))

	if err := createParentIfNecessary(destination); err != nil {
		log.Errorf("cannot create parent directories for path \"%s\": %s", path, err)
	}

	log.Infof("fetching image to %s", destination)
	if err := client.GetFileWithLogin(url, user, pass, destination); err != nil {
		log.Errorf("error fetching image: %s", err)
	}
}

// createParentIfNecessary takes the path provided and checks if its parent directory exists. If not, it creates it.
func createParentIfNecessary(path string) error {
	parent := filepath.Dir(path)

	// Check existence
	exists, err := pathExists(path)
	if err != nil {
		return fmt.Errorf("cannot check existence of path \"%s\": %w", parent, err)
	}

	// If exists, do nothing
	if exists {
		return nil
	}

	// Create parents if they don't exists
	log.Infof("creating parent dir: %s", parent)
	if err = os.MkdirAll(parent, 0755); err != nil {
		return fmt.Errorf("cannot create directory \"%s\": %w", parent, err)
	}
	return nil
}

// pathExists returns if the path provided exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
