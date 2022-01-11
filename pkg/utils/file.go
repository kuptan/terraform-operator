package utils

import (
	"os"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func CreateFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func WriteFile(filePath string, content string) {
	file, _ := CreateFile(filePath)

	file.WriteString(content)
}

func Mkdir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}
}
