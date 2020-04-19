package main

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

const (
	NGROK_EXECUTABLE_NAME = "ngrok"
	NGROK_ZIP_NAME        = "ngrok.zip"
)

func downloadNgrok() error {
	var ngrokUrl string
	armArch := regexp.MustCompile(`.*arm.*`)
	androidArch := regexp.MustCompile(`.*Android.*`)
	switch {
	case
		armArch.MatchString(runtime.GOARCH),
		androidArch.MatchString(runtime.GOARCH):
		ngrokUrl = "https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-arm.zip"
	default:
		ngrokUrl = "https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-386.zip"
	}

	wgetExec, _ := exec.LookPath("wget")

	wget := &exec.Cmd{
		Path:   wgetExec,
		Args:   []string{wgetExec, "--no-check-certificate", ngrokUrl, "-O", NGROK_ZIP_NAME},
		Stdout: nil,
		Stderr: nil,
	}

	if err := wget.Run(); err != nil {
		return err
	}
	return nil

}

func setNgrokAfterDownloading() error {
	unzipExec, _ := exec.LookPath("unzip")

	unzip := &exec.Cmd{
		Path:   unzipExec,
		Args:   []string{unzipExec, NGROK_ZIP_NAME},
		Stdout: nil,
		Stderr: os.Stdout,
	}

	if err := unzip.Run(); err != nil {
		return err
	}

	if err := os.Remove(NGROK_ZIP_NAME); err != nil {
		return err
	}

	chmodExec, _ := exec.LookPath("chmod")

	chmod := &exec.Cmd{
		Path:   chmodExec,
		Args:   []string{chmodExec, "+x", NGROK_EXECUTABLE_NAME},
		Stdout: nil,
		Stderr: os.Stdout,
	}

	if err := chmod.Run(); err != nil {
		return err
	}

	return nil
}

func isFileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return true, err
	}
}

func setNgrok() error {
	if check, err := isFileExists(NGROK_EXECUTABLE_NAME); check == true && err == nil {
		return nil
	} else if check == true && err != nil {
		return err
	} else {

		if err := downloadNgrok(); err != nil {
			return err
		}
		if err := setNgrokAfterDownloading(); err != nil {
			return err
		}

		return nil
	}
}
