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
