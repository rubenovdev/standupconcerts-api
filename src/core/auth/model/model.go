package model

import (
	"comedians/src/core/usersConcerts/model"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

type SignInDto struct {
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserDto struct {
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type PasswordRecoveryDto struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthGoogleDto struct {
	Email  string `json:"email" binding:"required,email"`
	ImgUrl string `json:"imgUrl"`
	Name   string `json:"name" binding:"required"`
	Id     string `json:"id" binding:"required"`
}

type AuthVkDto struct {
	Email  string `json:"email" binding:"required,email"`
	ImgUrl string `json:"imgUrl"`
	Name   string `json:"name" binding:"required"`
	Id     uint64 `json:"id" binding:"required"`
}

type AuthYandexDto struct {
	Email  string `json:"email" binding:"required,email"`
	ImgUrl string `json:"imgUrl"`
	Name   string `json:"name" binding:"required"`
	Id     string `json:"id" binding:"required"`
}

type Token struct {
	jwt.StandardClaims
	Id    uint64
	Roles []model.Role
}
