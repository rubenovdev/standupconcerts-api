package model

import (
	"time"
)

// cycle import fix

type User struct {
	Id                uint64     `json:"id" gorm:"primaryKey"`
	ImgUrl            *string    `json:"imgUrl,omitempty"`
	Password          string     `json:"-"`
	Email             string     `json:"email" gorm:"uniqueIndex"`
	Name              string     `json:"name"`
	About             string     `json:"about"`
	Concerts          *[]Concert `json:"concerts" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Subscriptions     []*Concert `json:"subscriptions" gorm:"many2many:users_subscriptions"`
	FavoriteConcerts  []*Concert `json:"favoriteConcerts" gorm:"many2many:users_favorite_concerts"`
	FavoriteComedians []*User    `json:"favoriteComedians" gorm:"many2many:users_favorte_comedians"`
	Roles             []Role     `json:"roles" gorm:"many2many:users_roles"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"default:current_timestamp"`
	UsersLikes        []*User    `json:"-" gorm:"many2many:users_comedians_likes"`
	UsersDislikes     []*User    `json:"-" gorm:"many2many:users_comedians_dislikes"`
	LikesCount        *uint64    `json:"likesCount" gorm:"default:0;not null"`
	DislikesCount     *uint64    `json:"dislikesCount" gorm:"default:0;not null"`
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
	Id          uint64    `json:"id" gorm:"primaryKey"`
	Filepath    string    `json:"filepath" gorm:"not null"`
	LikesCount  *uint64   `json:"likesCount" gorm:"default:0;not null"`
	ViewsCount  *int64    `json:"viewsCount" gorm:"default:0;not null"`
	UserId      uint64    `json:"userId" gorm:"not null"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title       string    `json:"title" gorm:"not null" binding:"required,min=2"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UsersLikes  []*User   `json:"-" gorm:"many2many:users_concerts_likes"`
}
