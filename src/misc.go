package main

import (
	"os"
	"runtime"
	"path/filepath"
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

func IsRegularFile(path string) (bool) {
	fileInfo, _ := os.Stat(path)
	return fileInfo.Mode().IsRegular()
}

func IsK8sConfigFile(path string) (bool) {
	if !IsRegularFile(path) {
		return false
	}

	switch(filepath.Ext(path)) {
	case ".json":
		return true
	case ".yaml":
		return true
	}

	return false
}
