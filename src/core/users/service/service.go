package service

import (
	commonService "comedians/src/common/service"
	concertsRepo "comedians/src/core/concerts/repo"
	"comedians/src/core/users/repo"
	"comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"
)

func UpdateUser(id uint64, user model.User) error {
	user.Id = id

	return repo.UpdateUser(user)
}

func DeleteUser(id uint64) error {
	return repo.DeleteUser(id)
}

func GetUser(id uint64) (model.User, error) {
	return repo.GetUser(id)
}

func GetUsers() ([]*model.User, error) {
	return repo.GetUsers()
}

func AppendFavoriteConcert(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}

	favoriteConcerts := append(user.FavoriteConcerts, &concert)

	return repo.UpdateFavoriteConcerts(user, favoriteConcerts)
}

func AppendFavoriteComedian(userId uint64, comedianId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	comedian, err := GetUser(comedianId)

	if err != nil {
		return err
	}

	favoriteComedians := append(user.FavoriteComedians, &comedian)

	return repo.UpdateFavoriteComedians(user, favoriteComedians)
}

func AppendSubscription(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}

	subscriptions := append(user.Subscriptions, &concert)

	return repo.UpdateSubscripions(user, subscriptions)
}

func DeleteFavoriteConcert(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}

	favoriteConcerts := utils.Filter(user.FavoriteConcerts, func(elem *model.Concert) bool {
		return elem.Id != concert.Id
	})

	return repo.UpdateFavoriteConcerts(user, favoriteConcerts)
}

func DeleteFavoriteComedian(userId uint64, comedianId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	comedian, err := GetUser(comedianId)

	if err != nil {
		return err
	}

	favoriteComedians := utils.Filter(user.FavoriteComedians, func(elem *model.User) bool {
		return elem.Id != comedian.Id
	})

	return repo.UpdateFavoriteComedians(user, favoriteComedians)
}

func UpdateUserImage(userId uint64, file multipart.File, filename string) (string, error) {
	var validExtensions = []string{".jpg", ".png"}

	if contains := commonService.ValidateExtension(filename, validExtensions); !contains {
		return "", errors.New("unsupported file")
	}

	dir := os.Getenv("USERS_IMAGES_DIR")
	commonService.MakeDirIfNotExists(dir)

	filepath := dir + "/" + fmt.Sprint(time.Now().Unix()) + ".jpg"
	user, _ := repo.GetUser(userId)

	if user.ImgUrl != nil {
		commonService.DeleteFile(*user.ImgUrl)
	}

	if err := commonService.UploadFile(file, filepath); err != nil {
		return "", err
	}

	user.ImgUrl = &filepath

	repo.UpdateUser(user)

	return filepath, nil
}

func DeleteSubscription(userId uint64, concertId uint64) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	concert, err := concertsRepo.GetConcert(concertId)

	if err != nil {
		return err
	}

	subscriptions := utils.Filter(user.Subscriptions, func(elem *model.Concert) bool {
		return elem.Id != concert.Id
	})

	return repo.UpdateSubscripions(user, subscriptions)
}

func UpdateUserPassword(userId uint64, password string) error {
	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	user.Password = utils.HashPassword(password)
	// todo
	err = UpdateUser(userId, user)

	return err
}

func Like(userId uint64, likingUserId uint64) error {
	likingUser, err := GetUser(likingUserId)

	if err != nil {
		return err
	}

	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	for _, userDislikes := range user.UsersDislikes {
		if userDislikes.Id == likingUserId {
			Dislike(userId, likingUserId)
			user, _ = GetUser(userId)
		}
	}

	contains := false

	for _, user := range user.UsersLikes {
		if user.Id == likingUserId {
			contains = true
		}
	}

	var usersLikes []*model.User

	if contains {
		usersLikes = utils.Filter(user.UsersLikes, func(user *model.User) bool {
			return user.Id != likingUserId
		})
	} else {
		usersLikes = append(user.UsersLikes, &likingUser)
	}

	user, err = UpdateUsersLikes(user, usersLikes)

	if contains {
		*user.LikesCount -= 1
	} else {
		*user.LikesCount += 1
	}

	if err != nil {
		return err
	}

	err = UpdateUser(user.Id, user)

	return err
}

func Dislike(userId uint64, dislikingUserId uint64) error {
	likingUser, err := GetUser(dislikingUserId)

	if err != nil {
		return err
	}

	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	for _, userLikes := range user.UsersLikes {
		if userLikes.Id == dislikingUserId {
			Like(user.Id, dislikingUserId)
			user, _ = GetUser(userId)
		}
	}

	contains := false

	for _, user := range user.UsersDislikes {
		if user.Id == dislikingUserId {
			contains = true
		}
	}

	var usersLikes []*model.User

	if contains {
		usersLikes = utils.Filter(user.UsersDislikes, func(user *model.User) bool {
			return user.Id != dislikingUserId
		})
	} else {
		usersLikes = append(user.UsersDislikes, &likingUser)
	}

	user, err = UpdateUsersDislikes(user, usersLikes)

	if contains {
		*user.DislikesCount -= 1
	} else {
		*user.DislikesCount += 1
	}

	if err != nil {
		return err
	}

	err = UpdateUser(user.Id, user)

	return err
}

func UpdateUsersLikes(user model.User, usersLikes []*model.User) (model.User, error) {
	return repo.UpdateUsersLikes(user, usersLikes)
}

func UpdateUsersDislikes(user model.User, usersDislikes []*model.User) (model.User, error) {
	return repo.UpdateUsersDislikes(user, usersDislikes)
}
