package controller

import (
	"comedians/src/common"
	authModel "comedians/src/core/auth/model"
	"comedians/src/core/auth/service"
	usersModel "comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

func signIn(c *gin.Context) {
	var input authModel.SignInDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user, err := service.GetUserByEmail(input.Email)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	signedToken, err := service.GenerateTokenJWT(user.Id, user.Roles)

	if err != nil {
		return
	}
	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"jwt": signedToken,
	}})
}

func signUp(c *gin.Context) {
	var input authModel.CreateUserDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	_, err := service.CreateUser(usersModel.User{Email: input.Email, Password: input.Password})

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func passwordRecovery(c *gin.Context) {
	var input authModel.PasswordRecoveryDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect data"))
		return
	}

	newPassword, err := service.RecoveryUserPassword(input.Email)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"password": newPassword,
	}})
}

func authGoogle(c *gin.Context) {
	var input authModel.AuthGoogleDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	signedToken, err := service.AuthGoogle(input)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"jwt": signedToken,
	}})
}

func authVk(c *gin.Context) {
	var input authModel.AuthVkDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	signedToken, err := service.AuthVk(input)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"jwt": signedToken,
	}})
}

func authYandex(c *gin.Context) {
	var input authModel.AuthYandexDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	signedToken, err := service.AuthYandex(input)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"jwt": signedToken,
	}})
}

func InitGroup(server *gin.Engine) *gin.RouterGroup {
	group := server.Group("auth")

	group.POST("/sign-in", signIn)
	group.POST("/sign-up", signUp)
	group.POST("/google", authGoogle)
	group.POST("/vk", authVk)
	group.POST("/yandex", authYandex)
	group.POST("/password-recovery", passwordRecovery)

	return group
}
