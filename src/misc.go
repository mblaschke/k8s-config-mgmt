package main

import (
	"os"
	"runtime"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func IsDirectory(path string) (bool) {
	fileInfo, _ := os.Stat(path)
	return fileInfo.IsDir()
}
