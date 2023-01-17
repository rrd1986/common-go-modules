package utils

import (
	"io/ioutil"
	"path/filepath"
)

type FileSystemType interface {
	ReadFile(filename string) ([]byte, error)
}

type FileSystem struct{}

func (fs FileSystem) ReadFile(filename string) ([]byte, error) {
	sanitisedFilePath := filepath.Clean(filename)
	file, err := ioutil.ReadFile(sanitisedFilePath)
	return file, err
}
