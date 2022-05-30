package utils

import (
	"os"
	"path/filepath"
)

// FileExists check if file exists in path
func FileExists(path string) bool {
	info, err := os.Stat(filepath.Clean(path))

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
