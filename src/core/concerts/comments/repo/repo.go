package repo

import (
	"comedians/src/core/concerts/comments/model"
	"comedians/src/db"

	"gorm.io/gorm"
)

const (
	commentsTable = "comments"
)

var commentsDB *gorm.DB

func lazyInit() {
	if commentsDB == nil {
		commentsDB = db.DBS.Table(commentsTable)
	}
}

func CreateComment(comment model.Comment) error {
	lazyInit()

	return commentsDB.Create(&comment).Error
}

func GetComments(concertId uint64) ([]model.Comment, error) {
	lazyInit()

	var comments []model.Comment

	err := commentsDB.Find(&comments, "concert_id = ?", concertId).Error

	return comments, err
}

func UpdateComment(comment model.Comment) error {
	lazyInit()

	return commentsDB.Where("id = ?", comment.Id).Updates(&comment).Error
}

func DeleteComment(concertId uint64, commentId uint64) error {
	lazyInit()

	var comments []model.Comment

	return commentsDB.Delete(&comments, "concert_id = ? & id = ?", concertId, commentId).Error
}
