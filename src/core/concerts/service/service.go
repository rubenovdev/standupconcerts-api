package service

import (
	commonService "comedians/src/common/service"
	"comedians/src/core/usersConcerts/model"	
	"comedians/src/core/concerts/repo"
	"comedians/src/utils"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

var concertsPath string

const (
	concertsDir = "concerts"
)

func init() {
	concertsPath = os.Getenv("ROOT_DIR") + "/" + concertsDir
	log.Print(concertsPath)
}

func GetConrerts() ([]model.Concert, error) {
	return repo.GetConcerts()
}

func GetConcert(id uint64) (model.Concert, error) {
	return repo.GetConcert(id)
} 

func UpdateConcert(concert model.Concert) (error) {
	return repo.UpdateConcert(concert)
} 

func CreateConcert(concert model.Concert) error {
	err := repo.CreateConcert(concert)
	return err
}

func UploadConcertFile(file multipart.File, filename string) (string, error) {
	var validExtensions = []string{".mp4", ".mov", ".wmv", ".avi", ".mpeg-4"}

	if contains := utils.Contains(validExtensions, filepath.Ext(filename)); !contains {
		return "", errors.New("unsupported file")
	}

	filepath := concertsPath + "/" + fmt.Sprint(time.Now().Unix()) + filename

	if err := commonService.UploadFile(file, filepath); err != nil {
		return "", err
	}

	return filepath, nil
}

func DeleteConcert(id uint64) (error) {
	return repo.DeleteConcert(id)
}

func DeleteConcertFile(filepath string) error {
	return commonService.DeleteFile(filepath)
}
