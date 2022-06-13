package service

import (
	"io"
	"mime/multipart"
	"os"
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
