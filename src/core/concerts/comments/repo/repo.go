package repo

import (
	commentsModel "comedians/src/core/concerts/comments/model"
	usersConcertsModel "comedians/src/core/usersConcerts/model"
	"comedians/src/db"
	"gorm.io/gorm"
)

const (
	commentsTable = "comments"
)

var commentsDB *gorm.DB

func lazyInit() {
	commentsDB = db.DBS.Preload("User").Preload("UsersLikes").Table(commentsTable)
}

func CreateComment(comment commentsModel.Comment) error {
	lazyInit()

	return commentsDB.Create(&comment).Error
}

func GetComments(concertId uint64) ([]commentsModel.Comment, error) {
	lazyInit()

	var comments []commentsModel.Comment

	err := commentsDB.Order("created_at").Find(&comments, "concert_id = ?", concertId).Error

	return comments, err
}

func UpdateComment(comment commentsModel.Comment) error {
	lazyInit()

	return commentsDB.Where("id = ?", comment.Id).Updates(&comment).Error
}

func DeleteComment(concertId uint64, commentId uint64) error {
	lazyInit()

	var comments []commentsModel.Comment

	return commentsDB.Delete(&comments, "concert_id = ? AND id = ?", concertId, commentId).Error
}

func GetComment(concertId uint64, commentId uint64) (commentsModel.Comment, error) {
	lazyInit()

	var comment commentsModel.Comment

	err := commentsDB.First(&comment, "concert_id = ? AND id = ?", concertId, commentId).Error

	return comment, err
}

func UpdateUsersLikes(comment commentsModel.Comment, usersLikes []*usersConcertsModel.User) (commentsModel.Comment, error) {
	lazyInit()

	err := db.DBS.Model(&comment).Association("UsersLikes").Clear()

	if err != nil {
		return commentsModel.Comment{}, err
	}

	err = db.DBS.Save(&comment).Error

	if err != nil {
		return commentsModel.Comment{}, err
	}

	comment.UsersLikes = usersLikes

	err = db.DBS.Save(&comment).Error
	return comment, err
}
