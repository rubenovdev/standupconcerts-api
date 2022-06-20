package controller

import (
	"comedians/src/common"
	authModel "comedians/src/core/auth/model"
	"comedians/src/core/auth/service"
	usersModel "comedians/src/core/usersConcerts/model"
	"comedians/src/utils"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signIn(c *gin.Context) {
	var input authModel.SignInDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect data"))
		return
	}

	log.Print("input sign in ", input)

	user, err := service.GetUserByEmail(input.Email)

	log.Print(user)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect data"))
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect data"))
		return
	}

	signedToken, err := service.GenerateTokenJWT(user.Id, user.Roles)

	if err != nil {
		log.Panic(err)
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

	log.Print("input", input.Email)

	err := service.CreateUser(usersModel.User{Email: input.Email, Password: input.Password})

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

func InitGroup(server *gin.Engine) *gin.RouterGroup {
	group := server.Group("auth")

	group.POST("/sign-in", signIn)
	group.POST("/sign-up", signUp)
	group.POST("/password-recovery", passwordRecovery)

	return group
}
