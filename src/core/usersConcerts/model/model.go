package model

import (
	"time"
)

// cycle import fix

type User struct {
	Id                uint64    `json:"id" gorm:"primaryKey"`
	ImgUrl            *string   `json:"imgUrl,omitempty"`
	Password          string    `json:"-"`
	Email             string    `json:"email" gorm:"uniqueIndex"`
	Name              string    `json:"name,omitempty"`
	About             string    `json:"about,omitempty"`
	Subscriptions     []*Concert `json:"subscriptions" gorm:"many2many:users_subscriptions"`
	FavoriteConcerts  []*Concert `json:"favoriteConcerts" gorm:"many2many:users_favorite_concerts"`
	FavoriteComedians []*User    `json:"favoriteComedians" gorm:"many2many:users_favorte_comedians"`
	Roles             []Role    `json:"roles" gorm:"many2many:users_roles"`
	CreatedAt         time.Time `json:"created_at" gorm:"default:current_timestamp"`
}

type Role struct {
	Id          uint64       `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title"`
	Permissions []Permission `json:"permissions" gorm:"many2many:roles_permissions;"`
}

type Permission struct {
	Id    uint64 `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}

type Concert struct {
	Id         uint64    `json:"id" gorm:"primaryKey"`
	Filepath   string    `json:"filepath" gorm:"not null"`
	LikeCount  int64     `json:"likeCount" gorm:"default:0;not null"`
	WatchCount int64     `json:"watchCount" gorm:"default:0;not null"`
	UserId     uint64    `json:"userId" gorm:"not null"`
	User       User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title      string    `json:"title" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:current_timestamp"`
}
