package sio

import (
	"os"

	"github.com/syncfuture/go/serr"
)

func Exists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	info.IsDir()

	if os.IsExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	// 其他错误
	return false, serr.WithStack(err)
}

func GetDir(path string) (string, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return path, nil
		}
		return "", nil
	}

	if info.IsDir() {
		return path, nil
	}

	if os.IsExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	// 其他错误
	return false, serr.WithStack(err)
}

func Write(path string, data []byte) error {
	exists, err := Exists(path)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 666)
	return serr.WithStack(err)
}
