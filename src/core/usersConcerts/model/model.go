package model

import (
	"time"
)

// cycle import fix

type User struct {
	Id                uint64     `json:"id" gorm:"primaryKey"`
	ImgUrl            *string    `json:"imgUrl,omitempty"`
	Password          string     `json:"-"`
	Email             string     `json:"email" gorm:"uniqueIndex" binding:"email"`
	Name              string     `json:"name"`
	About             string     `json:"about"`
	Concerts          *[]Concert `json:"concerts" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Subscriptions     []*Concert `json:"subscriptions" gorm:"many2many:users_subscriptions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FavoriteConcerts  []*Concert `json:"favoriteConcerts" gorm:"many2many:users_favorite_concerts;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FavoriteComedians []*User    `json:"favoriteComedians" gorm:"many2many:users_favorte_comedians;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Roles             []Role     `json:"roles" gorm:"many2many:users_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"default:current_timestamp"`
	UsersLikes        []*User    `json:"-" gorm:"many2many:users_comedians_likes;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UsersDislikes     []*User    `json:"-" gorm:"many2many:users_comedians_dislikes;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LikesCount        *uint64    `json:"likesCount" gorm:"default:0;not null"`
	DislikesCount     *uint64    `json:"dislikesCount" gorm:"default:0;not null"`
}

type Role struct {
	Id          uint64       `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title"`
	Permissions []Permission `json:"permissions" gorm:"many2many:roles_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Permission struct {
	Id    uint64 `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}

type Concert struct {
	Id          uint64    `json:"id" gorm:"primaryKey"`
	VideoSrc    string    `json:"videoSrc" gorm:"not null"`
	FrameSrc    string    `json:"frameSrc" gorm:"not null"`
	LikesCount  *uint64   `json:"likesCount" gorm:"default:0;not null"`
	ViewsCount  *int64    `json:"viewsCount" gorm:"default:0;not null"`
	UserId      uint64    `json:"userId" gorm:"not null"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title       string    `json:"title" gorm:"not null" binding:"required,min=2"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UsersLikes  []*User   `json:"-" gorm:"many2many:users_concerts_likes;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
