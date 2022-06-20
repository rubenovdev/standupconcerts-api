package service

import (
	commonService "comedians/src/common/service"
	concertsModel "comedians/src/core/concerts/model"
	"comedians/src/core/concerts/repo"
	usersService "comedians/src/core/users/service"
	usersConcertsModel "comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
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

func GetConrerts(filters concertsModel.Filters) ([]usersConcertsModel.Concert, error) {
	if filters.SortBy == "new" {
		filters.SortBy = "created_at"
	} else if filters.SortBy == "popular" {
		filters.SortBy = "views_count"
	} else {
		filters.SortBy = "created_at"
	}

	log.Print(filters.ComedianId, filters.SortBy, filters.Year)
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

func UploadConcertFile(file multipart.File, filename string) (string, error) {
	var validExtensions = []string{".mp4", ".mov", ".wmv", ".avi", ".mpeg-4"}

	if contains := commonService.ValidateExtension(filename, validExtensions); !contains {
		return "", errors.New("unsupported file")
	}
	filepath := concertsPath + "/" + fmt.Sprint(time.Now().Unix()) + ".mp4"

	if err := commonService.UploadFile(file, filepath); err != nil {
		return "", err
	}

	return filepath, nil
}

func DeleteConcert(id uint64) error {
	concert, err := GetConcert(id)

	if err != nil {
		return err
	}

	commonService.DeleteFile(concert.Filepath)

	return repo.DeleteConcert(id)
}

func DeleteConcertFile(filepath string) error {
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
		log.Print("userId", user.Id, "Second", userId)

		if user.Id == userId {
			contains = true
		}
	}

	log.Print(contains)

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
