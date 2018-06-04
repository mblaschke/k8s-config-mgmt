package config

import (
	"os"
	"path/filepath"
)

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

func recursiveFileListByPath(path string) (list []string) {
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if IsK8sConfigFile(path) {
			list = append(list, path)
		}
		return nil
	})

	return
}

func ensureAbsConfigPath(basePath, path string) (absPath string) {
	var err error
	absPath = path

	if !filepath.IsAbs(absPath) {
		absPath, err = filepath.Abs(filepath.Join(filepath.Dir(basePath), absPath))
		if err != nil {
			panic(err.Error())
		}
	}

	return
}
