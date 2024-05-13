package utils

import (
	"archive/tar"
	"bytes"
	"os"
	"path"
	"strings"
)

type File struct {
	Name  string `json:"name"`
	Input string `json:"input"`
}

type CreateTarOptions struct {
	FileMode int
	FileUid  int
	FileGid  int
	DirMode  int
	DirUid   int
	DirGid   int
}

func CreateTar(files []File, defaultOptions *CreateTarOptions) (*bytes.Buffer, error) {
	if defaultOptions == nil {
		defaultOptions = &CreateTarOptions{
			FileMode: 0655,
			FileUid:  1000,
			FileGid:  1001,
			DirMode:  0755,
			DirUid:   1000,
			DirGid:   1001,
		}
	}

	// Create a buffer to hold the archive data
	var buffer bytes.Buffer

	// Create a new tar archive writer
	tw := tar.NewWriter(&buffer)
	defer tw.Close()

	// Loop through each file
	for _, file := range files {
		// Split the file path into components
		filePathParts := strings.Split(file.Name, string(os.PathSeparator))

		// Iterate through path components except the filename
		for i := range filePathParts[:len(filePathParts)-1] {
			// Create a directory entry for each parent folder
			parentDir := path.Join(filePathParts[:i+1]...)
			dirHeader := new(tar.Header)
			dirHeader.Typeflag = tar.TypeDir
			dirHeader.Name = parentDir
			dirHeader.Mode = int64(defaultOptions.DirMode)
			dirHeader.Uid = defaultOptions.DirUid
			dirHeader.Gid = defaultOptions.DirGid

			// Write the directory entry to the archive
			if err := tw.WriteHeader(dirHeader); err != nil {
				return nil, err
			}
		}

		// Create a tar header for the file
		header := new(tar.Header)
		header.Name = file.Name
		header.Size = int64(len(file.Input))
		header.Mode = int64(defaultOptions.FileMode)
		header.Uid = defaultOptions.FileUid
		header.Gid = defaultOptions.FileGid

		// Write the file content to the archive
		if err := tw.WriteHeader(header); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(file.Input)); err != nil {
			return nil, err
		}
	}

	return &buffer, nil
}
