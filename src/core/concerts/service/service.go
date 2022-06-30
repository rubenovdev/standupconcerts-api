package service

import (
	"bytes"
	commonService "comedians/src/common/service"
	"comedians/src/core/concerts/model"
	concertsModel "comedians/src/core/concerts/model"
	"comedians/src/core/concerts/repo"
	usersService "comedians/src/core/users/service"
	usersConcertsModel "comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func GetConrerts(filters concertsModel.Filters) ([]usersConcertsModel.Concert, error) {
	if filters.SortBy == "new" {
		filters.SortBy = "created_at"
	} else if filters.SortBy == "popular" {
		filters.SortBy = "views_count"
	} else {
		filters.SortBy = "created_at"
	}

	return repo.GetConcerts(filters)
}

func GetConcert(id uint64) (usersConcertsModel.Concert, error) {
	return repo.GetConcert(id)
}

func UpdateConcert(concert usersConcertsModel.Concert) error {
	return repo.UpdateConcert(concert)
}

func CreateConcert(concert usersConcertsModel.Concert) error {
	err := repo.CreateConcert(concert)
	return err
}

func UploadConcertFile(file multipart.File, primalFilename string) (string, string, error) {
	var validExtensions = []string{".mp4", ".mov", ".wmv", ".avi", ".mpeg-4"}

	if contains := commonService.ValidateExtension(primalFilename, validExtensions); !contains {
		return "", "", errors.New("unsupported file")
	}

	concertsFramesPath := os.Getenv("CONCERTS_FRAMES_DIR")
	concertsVideosPath := os.Getenv("CONCERTS_VIDEOS_DIR")

	commonService.MakeDirIfNotExists(concertsFramesPath)
	commonService.MakeDirIfNotExists(concertsVideosPath)

	filename := fmt.Sprint(time.Now().Unix()) + ""

	filepathVideo := concertsVideosPath + "/" + filename + ".mp4"
	filepathFrame := concertsFramesPath + "/" + filename + ".jpg"

	if err := commonService.UploadFile(file, filepathVideo); err != nil {
		return "", "", err
	}

	err := commonService.ExtractFrames(filepathVideo, filepathFrame, 1)

	if err != nil {
		return "", "", err
	}

	return filepathVideo, filepathFrame, nil
}

func DownloadVideoFromYoutube(link string) (string, string, error) {
	concertsFramesPath := os.Getenv("CONCERTS_FRAMES_DIR")
	concertsVideosPath := os.Getenv("CONCERTS_VIDEOS_DIR")

	commonService.MakeDirIfNotExists(concertsFramesPath)
	commonService.MakeDirIfNotExists(concertsVideosPath)

	filename := fmt.Sprint(time.Now().Unix())

	body := new(model.YoutubeDownloadBody)
	body.Link = link
	body.Filename = filename + ".mp4"
	body.OutDir = commonService.GetRootDir() + "/" + concertsVideosPath

	videoFilepath := concertsVideosPath + "/" + filename + ".mp4"
	frameFilepath := concertsFramesPath + "/" + filename + ".jpg"

	jsonBytes, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", os.Getenv("YOUTUBE_SERVICE_API")+"/upload-video", bytes.NewBuffer(jsonBytes))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		return "", "", err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return "", "", err
	}

	err = commonService.ExtractFrames(videoFilepath, frameFilepath, 1)

	if err != nil {
		return "", "", err
	}

	return videoFilepath, frameFilepath, nil
}

func DeleteConcert(id uint64) error {
	concert, err := GetConcert(id)

	if err != nil {
		return err
	}

	commonService.DeleteFile(concert.VideoSrc)
	commonService.DeleteFile(concert.FrameSrc)

	return repo.DeleteConcert(id)
}

func DeleteConcertVideo(filepath string) error {
	return commonService.DeleteFile(filepath)
}

func DeleteConcertFrame(filepath string) error {
	return commonService.DeleteFile(filepath)
}

func Like(concertId uint64, userId uint64) error {
	concert, err := GetConcert(concertId)
	user, _ := usersService.GetUser(userId)

	if err != nil {
		return err
	}

	contains := false

	for _, user := range concert.UsersLikes {

		if user.Id == userId {
			contains = true
		}
	}

	var usersLikes []*usersConcertsModel.User

	if contains {
		usersLikes = utils.Filter(concert.UsersLikes, func(user *usersConcertsModel.User) bool {
			return user.Id != userId
		})
	} else {
		usersLikes = append(concert.UsersLikes, &user)
	}

	concert, err = UpdateUsersLikes(concert, usersLikes)

	if contains {
		*concert.LikesCount -= 1
	} else {
		*concert.LikesCount += 1
	}

	if err != nil {
		return err
	}

	err = UpdateConcert(concert)

	return err
}

func UpdateUsersLikes(concert usersConcertsModel.Concert, usersLikes []*usersConcertsModel.User) (usersConcertsModel.Concert, error) {
	return repo.UpdateUsersLikes(concert, usersLikes)
}

func IncViews(concertId uint64) error {
	concert, err := GetConcert(concertId)

	if err != nil {
		return err
	}

	*concert.ViewsCount += 1
	err = UpdateConcert(concert)

	return err
}
