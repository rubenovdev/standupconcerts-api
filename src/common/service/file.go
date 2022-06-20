package service

import (
	"comedians/src/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadFile(file multipart.File, filepath string) error {
	defer file.Close()

	// Create file
	dst, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

func MakeDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func ValidateExtension(filename string, admissibleExtensions []string) bool {
	return utils.Contains(admissibleExtensions, filepath.Ext(filename))
}
