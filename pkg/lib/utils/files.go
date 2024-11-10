package utils

import (
	"fmt"
	"os"
)

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func MakeDirIfNotExist(src string) error {
	if CheckNotExist(src) {
		return os.MkdirAll(src, os.ModePerm)
	}
	return nil
}

func MustOpen(dir, fileName string) (*os.File, error) {
	if CheckPermission(dir) {
		return nil, fmt.Errorf("permission denied")
	}
	if err := MakeDirIfNotExist(dir); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
