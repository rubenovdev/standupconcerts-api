package service

import (
	"comedians/src/core/concerts/comments/model"
	"comedians/src/core/concerts/comments/repo"
)

func CreateComment(comment model.Comment) (error) {
	return repo.CreateComment(comment)
}

func GetComments(concertId uint64) ([]model.Comment, error) {
	return repo.GetComments(concertId)
}

func DeleteComment(concertId uint64, commentId uint64) (error) {
	return repo.DeleteComment(concertId, commentId)
}

func UpdateComment(comment model.Comment) (error) {
	return repo.UpdateComment(comment)
}