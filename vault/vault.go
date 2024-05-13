package main

import (
	"fmt"
	"molta9a/server/utils"
	"os"
	"path/filepath"
	"strings"
)

var VAULT_PATH = os.Getenv("VAULT_PATH")

func main() {
	res, err := ReadFile("courses/javascript/TAGS.md")
	fmt.Println(string(res))
	fmt.Println("-----------------------------------------------")
	fmt.Println(err)
}

func Read(relativePath string) (*[]byte, *[]string, error) {
	realPath, err := cleanPath(relativePath)
	if err != nil {
		return nil, nil, err
	}

	fileInfo, err := os.Stat(realPath)
	if err != nil {
		return nil, nil, err
	}

	if fileInfo.IsDir() {
		files, err := unsafeReadDirectory(realPath)
		if err != nil {
			return nil, nil, err
		}
		return nil, &files, nil
	}
	// else
	data, err := unsafeReadFile(realPath)
	if err != nil {
		return nil, nil, err
	}
	return &data, nil, err
}

func ReadFile(relativePath string) ([]byte, error) {
	realPath, err := cleanPath(relativePath)
	if err != nil {
		return nil, err
	}
	return unsafeReadFile(realPath)
}

func ReadDirectory(relativePath string) ([]string, error) {
	realPath, err := cleanPath(relativePath)
	if err != nil {
		return nil, err
	}
	return unsafeReadDirectory(realPath)
}

func cleanPath(relativePath string) (string, error) {
	fullPath := filepath.Join(VAULT_PATH, relativePath)

	// Ensure that the full path is within the allowed directory
	if !strings.HasPrefix(fullPath, VAULT_PATH) {
		return "", utils.ErrAccessDenied
	}

	// Sanitize the file path to remove ".." and resolve symbolic links
	cleanPath := filepath.Clean(fullPath)
	cleanPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		return "", utils.ErrAccessDenied
	}

	// Ensure that the sanitized path is still within the allowed directory
	if !strings.HasPrefix(cleanPath, VAULT_PATH) {
		return "", utils.ErrAccessDenied
	}

	return cleanPath, nil
}

func unsafeReadFile(path string) ([]byte, error) {
	realPath, err := cleanPath(path)
	if err != nil {
		return nil, err
	}

	// Read file
	data, err := os.ReadFile(realPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func unsafeReadDirectory(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}
