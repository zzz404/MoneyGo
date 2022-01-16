package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

func Must(errs ...error) {
	var sb strings.Builder
	errCount := 0
	for _, err := range errs {
		if err != nil {
			if errCount > 0 {
				sb.WriteString(";\n")
			}
			sb.WriteString(err.Error())
			errCount++
		}
	}
	if errCount > 0 {
		panic(errors.New(sb.String()))
	}
}

func FormatDate(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02")
	} else {
		return ""
	}
}

func FormatTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	} else {
		return ""
	}
}

func ParseDate(s string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	return &t, err
}

func ParseTime(s string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	return &t, err
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
