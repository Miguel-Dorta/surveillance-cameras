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
	ip, user, pass, camName, path               string
	rtspPort                                    int
	recordingDuration, endingTimeout, retryTime time.Duration

	ffmpegPath string
	verbose bool
	log     *logolang.Logger
	cmd     *command
)

func init() {
	log = logolang.NewLogger()
	log.Color = false
	log.Level = logolang.LevelError

	var (
		pidDir  string
		err error
		version bool
	)
	flag.StringVar(&ip, "ip", "", "Camera IP")
	flag.IntVar(&rtspPort, "rtsp-port", 554, "Camera RTSP Port")
	flag.StringVar(&user, "user", "", "Username for login")
	flag.StringVar(&pass, "password", "", "Password for login")
	flag.StringVar(&camName, "camera-name", "", "Camera ID")
	flag.StringVar(&path, "path", "", "Path to save")
	flag.StringVar(&pidDir, "pid-directory", "/run", "Path to pid file's directory")
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
	case ip:
		logCritical("invalid ip")
	case user:
		logCritical("invalid username")
	case pass:
		logCritical("invalid password")
	case camName:
		logCritical("invalid camera-name")
	case path:
		logCritical("invalid path destination")
	case pidDir:
		logCritical("invalid pid path")
	}

	if rtspPort < 0 || rtspPort > 65535 {
		logCritical("invalid port")
	}

	if err := si.Register("OWIPPB200_downloadVideo-" + camName); err != nil {
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
	now := time.Now().UTC()
	ctx, _ := context.WithTimeout(context.Background(), recordingDuration+endingTimeout)
	execCmd := exec.CommandContext(ctx,
		ffmpegPath, "-rtsp_transport", "tcp", "-i",
		fmt.Sprintf("rtsp://%s:%s@%s:%d/11", user, pass, ip, rtspPort),
		filepath.Join(
			path,
			camName,
			strconv.Itoa(now.Year()),
			strconv.Itoa(int(now.Month())),
			strconv.Itoa(now.Day()),
			fmt.Sprintf("%02d-%02d-%02d.mp4", now.Hour(), now.Minute(), now.Second())))
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
