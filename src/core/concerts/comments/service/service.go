package service

import (
	commentsModel "comedians/src/core/concerts/comments/model"
	"comedians/src/core/concerts/comments/repo"
	usersService "comedians/src/core/users/service"
	usersConcertsModel "comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
)

func CreateComment(comment commentsModel.Comment) error {
	return repo.CreateComment(comment)
}

func GetComments(concertId uint64) ([]commentsModel.Comment, error) {
	return repo.GetComments(concertId)
}

func DeleteComment(concertId uint64, commentId uint64) error {
	return repo.DeleteComment(concertId, commentId)
}

func UpdateComment(comment commentsModel.Comment) error {
	return repo.UpdateComment(comment)
}

func GetComment(concertId uint64, commentId uint64) (commentsModel.Comment, error) {
	return repo.GetComment(concertId, commentId)
}

func Like(concertId uint64, commentId uint64, userId uint64) error {
	comment, err := GetComment(concertId, commentId)
	user, _ := usersService.GetUser(userId)

	if err != nil {
		return err
	}

	contains := false

	for _, user := range comment.UsersLikes {
		if user.Id == userId {
			contains = true
		}
	}

	var usersLikes []*usersConcertsModel.User

	if contains {
		usersLikes = utils.Filter(comment.UsersLikes, func(user *usersConcertsModel.User) bool {
			return user.Id != userId
		})
	} else {
		usersLikes = append(comment.UsersLikes, &user)
	}

	comment, err = UpdateUsersLikes(comment, usersLikes)

	if err != nil {
		return err
	}

	if contains {
		*comment.LikesCount -= 1
	} else {
		*comment.LikesCount += 1
	}

	err = UpdateComment(comment)

	return err
}

func UpdateUsersLikes(comment commentsModel.Comment, usersLikes []*usersConcertsModel.User) (commentsModel.Comment, error) {
	return repo.UpdateUsersLikes(comment, usersLikes)
}
