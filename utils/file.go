package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

const BUF_SIZE = 1024

// 文件写入内容
func FilePutContent(file string, content string) error {
	// 创建目录
	if _, err := os.Stat(file); err != nil {
		if err = os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			return err
		}
		if _, err = os.Create(file); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(file, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return err
	}
	if n, err := f.WriteString(content); err != nil {
		return err
	} else if n != len(content) {
		return errors.New("写入数据不全")
	}

	return nil
}

// 文件读出内容
func FileGetContent(file string) (string, error) {
	content := ""
	if _, err := os.Stat(file); err != nil {
		return content, nil
	}
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return content, err
	}
	defer f.Close()
	contentBytes := make([]byte, BUF_SIZE)
	for true {
		if n, err := f.Read(contentBytes); err != nil {
			if err == io.EOF {
				return content, nil
			}
			return "", err
		} else {
			content += string(contentBytes[:n])
			if n < BUF_SIZE {
				break
			}
		}
	}
	return content, nil
}
