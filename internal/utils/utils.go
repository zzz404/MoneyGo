package utils

import (
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

func CopyFile(srcPath, destPath string) error {
	if err := AssertFileExists(srcPath); err != nil {
		return err
	}
	if err := AssertFileNotExists(destPath); err != nil {
		return err
	}
	source, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
