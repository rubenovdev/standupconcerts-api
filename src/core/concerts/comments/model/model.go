package model

import (
	usersConcertsModel "comedians/src/core/usersConcerts/model"
)

type Comment struct {
	Id        uint64                     `json:"id" gorm:"primaryKey"`
	UserId    uint64                     `json:"userId" gorm:"not null"`
	User      usersConcertsModel.User    `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ConcertId uint64                     `json:"concertId" gorm:"not null"`
	Concert   usersConcertsModel.Concert `json:"concert" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Body      string                     `json:"body" binding:"required"`
	LikeCount int64                      `json:"likeCount" gorm:"default:0;not null"`
}

type CreateCommentDto struct {
	Body   string `json:"body" validate:"presence,min=1,max=200"`
	UserId uint64 `json:"userId"`
}
