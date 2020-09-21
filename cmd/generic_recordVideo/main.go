package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/si"
	"github.com/Miguel-Dorta/surveillance-cameras/internal"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type command struct {
	*exec.Cmd
	stdinPipe io.WriteCloser
}

var (
	url, camName, path                          string
	recordingDuration, endingTimeout, retryTime time.Duration

	ffmpegPath string
	verbose    bool
	log        *logolang.Logger
	cmd        *command
)

func init() {
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelError

	var (
		err     error
		version bool
	)
	flag.StringVar(&url, "url", "", "RTSP stream URL (example: rtsp://user:pass@127.0.0.1:554/myStream")
	flag.StringVar(&camName, "camera-name", "", "Camera ID")
	flag.StringVar(&path, "path", "", "Path to save")
	flag.StringVar(&si.Dir, "pid-directory", "/run", "Path to pid file's directory")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&version, "V", false, "Print version and exit")
	flag.DurationVar(&recordingDuration, "recording-duration", time.Minute*10, "Set the duration of each clip")
	flag.DurationVar(&endingTimeout, "ending-timeout", time.Minute, "Set the timeout before killing ffmpeg. A killed process will leave a corrupt file.")
	flag.DurationVar(&retryTime, "retry-every", time.Minute*10, "Set the time between retries for executing ffmpeg.")
	flag.Parse()

	if version {
		fmt.Println(internal.Version)
		os.Exit(0)
	}

	ffmpegPath, err = exec.LookPath("ffmpeg")
	if err != nil {
		logCritical("dependency ffmpeg not found")
	}

	// Check correct args
	switch "" {
	case url:
		logCritical("invalid URL")
	case camName:
		logCritical("invalid camera-name")
	case path:
		logCritical("invalid path destination")
	case si.Dir:
		logCritical("invalid pid path")
	}

	if err := si.Register("generic_downloadVideo_" + camName); err != nil {
		if err == si.ErrOtherInstanceRunning {
			os.Exit(0)
		} else {
			log.Criticalf("error registering program instance: %s", err)
			os.Exit(1)
		}
	}
}

func main() {
	for {
		if err := start(); err != nil {
			log.Error(err.Error())
			time.Sleep(retryTime)
			continue
		}

		time.Sleep(recordingDuration)

		if err := stop(); err != nil {
			log.Error(err.Error())
		}
	}
}

func start() error {
	dir, filename := getNewFilePath()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("error creating parent directory: %s", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), recordingDuration+endingTimeout)
	execCmd := exec.CommandContext(ctx, ffmpegPath, "-rtsp_transport", "tcp", "-i", url, "-acodec", "copy", "-vcodec", "copy", filepath.Join(dir, filename))
	if verbose {
		execCmd.Stderr = os.Stderr
		execCmd.Stdout = os.Stdout
	}

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %w", err)
	}

	if err := execCmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	cmd = &command{
		Cmd:       execCmd,
		stdinPipe: stdin,
	}

	return nil
}

func getNewFilePath() (string, string) {
	now := time.Now().UTC()
	return filepath.Join(
			path,
			camName,
			strconv.Itoa(now.Year()),
			fmt.Sprintf("%02d", now.Month()),
			fmt.Sprintf("%02d", now.Day())),
		fmt.Sprintf("%02d-%02d-%02d.mkv", now.Hour(), now.Minute(), now.Second())
}

func stop() error {
	if _, err := cmd.stdinPipe.Write([]byte{'q'}); err != nil {
		return fmt.Errorf("error sending quit key: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error returning waiting to stop: %w", err)
	}
	return nil
}

func logCritical(msg string) {
	log.Critical(msg)
	os.Exit(1)
}
