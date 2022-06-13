package model

import (
	"comedians/src/core/usersConcerts/model"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

type SignInDto struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type CreateUserDto struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type PasswordRecoveryDto struct {
	Email string `json:"email" binding:"required"`
}

type Token struct {
	jwt.StandardClaims
	Id    uint64
	Roles []model.Role
}
