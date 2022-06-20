package model

type UpdateUserPasswordDto struct {
	Password string `json:"password" binding:"required,min=6"`
}
