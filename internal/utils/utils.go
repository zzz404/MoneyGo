package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func GetTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func AssertFileExists(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("檔案 %s 不存在", path)
		} else {
			return err
		}
	}
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("%s 不是正常檔案", path)
	}
	return nil
}

func AssertDirExists(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("目錄 %s 不存在", path)
		} else {
			return err
		}
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s 不是目錄", path)
	}
	return nil
}

func AssertFileNotExists(path string) error {
	if _, err := os.Stat(path); err == nil || os.IsExist(err) {
		return fmt.Errorf("%s 已存在", path)
	}
	return nil
}

func CopyFile(srcPath, destPath string) (err error) {
	if err = AssertFileExists(srcPath); err != nil {
		return
	}
	if err = AssertFileNotExists(destPath); err != nil {
		return
	}
	source, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer func() {
		err = CombineError(err, source.Close())
	}()

	destination, err := os.Create(destPath)
	if err != nil {
		return
	}
	defer func() {
		err = CombineError(err, destination.Close())
	}()

	_, err = io.Copy(destination, source)
	return
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func CombineError(err1 error, err2 error) error {
	if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return errors.New(err1.Error() + ";\n" + err2.Error())
	}
}
