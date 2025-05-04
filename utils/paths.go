package utils

import (
	"os"
	"path/filepath"
)

func RelPathFromCwd(path string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Rel(cwd, abs)
}
