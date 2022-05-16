package logger

import (
	"fmt"
	"os"
)

func mustOpen(fileName string, dir string) (*os.File, error) {
	if checkPermission(dir) {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}
	if err := isNotExistMkDir(dir); err != nil {
		return nil, fmt.Errorf("error during make dir: %s, err: %+v", dir, err)
	}
	file, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}
	return file, nil
}

func checkPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func isNotExistMkDir(src string) error {
	if notExist := checkNotExist(src); notExist {
		if err := mkdir(src); err != nil {
			return err
		}
	}
	return nil
}

func mkdir(src string) error {
	if err := os.MkdirAll(src, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}
