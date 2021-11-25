package localpath

import (
	"os"
	"path/filepath"
	"regexp"
)

func getExePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath := filepath.Dir(exe)
	return exePath, nil
}

func getWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

//TODO: Better handle poorly written regex expression errors
func Get() (string, error) {
	exePath, err := getExePath()

	if err != nil {
		return "", err
	}

	match, _ := regexp.MatchString("^(/tmp)", exePath)

	if !match {
		match, _ = regexp.MatchString("(/go-build)", exePath)
	}

	if !match {
		return exePath, nil
	}
	return getWorkingDirectory()
}
