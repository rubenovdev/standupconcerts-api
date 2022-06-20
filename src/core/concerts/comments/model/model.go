package model

import (
	usersConcertsModel "comedians/src/core/usersConcerts/model"
	"time"
)

type Comment struct {
	Id         uint64                     `json:"id" gorm:"primaryKey"`
	UserId     uint64                     `json:"userId" gorm:"not null"`
	User       usersConcertsModel.User    `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ConcertId  uint64                     `json:"concertId" gorm:"not null"`
	Concert    usersConcertsModel.Concert `json:"concert" gorm:"foreignKey:ConcertId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Body       string                     `json:"body" binding:"required,min=1"`
	LikesCount *uint64                     `json:"likesCount" gorm:"default:0;not null"`
	UsersLikes []*usersConcertsModel.User `json:"-" gorm:"many2many:users_comments_likes"`
	CreatedAt  time.Time                  `json:"createdAt" gorm:"default:current_timestamp"`
}

type CreateCommentDto struct {
	Body string `json:"body" binding:"required,min=1"`
}
