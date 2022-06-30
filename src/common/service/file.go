package service

import (
	"bytes"
	"comedians/src/utils"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

func ExtractFrames(inFilepath, outFilePath string, fps int) error {

	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFilepath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", fps)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		return err
	}

	img, err := imaging.Decode(buf)

	if err != nil {
		return err
	}

	err = imaging.Save(img, outFilePath)

	if err != nil {
		return err
	}
	return nil
}

func DownloadFile(url string, outFilePath string) error {
	out, err := os.Create(outFilePath)

	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func GetRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}